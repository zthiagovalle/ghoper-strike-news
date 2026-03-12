package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/zthiagovalle/ghoper-strike-news/internal/ai"
	"github.com/zthiagovalle/ghoper-strike-news/internal/config"
	"github.com/zthiagovalle/ghoper-strike-news/internal/database"
	"github.com/zthiagovalle/ghoper-strike-news/internal/discord"
	"github.com/zthiagovalle/ghoper-strike-news/internal/feed"
)

const (
	rssFeedURL    = "https://store.steampowered.com/feeds/news/app/730/"
	checkInterval = 15 * time.Minute
)

type Bot struct {
	repo     *database.Repository
	parser   *feed.RSSParser
	notifier *discord.Notifier
	analyzer *ai.Analyzer
}

func NewBot(cfg *config.Config) (*Bot, error) {
	repo, err := database.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	notifier, err := discord.NewNotifier(cfg.DiscordToken, cfg.ChannelID)
	if err != nil {
		repo.Close()
		return nil, err
	}

	var analyzer *ai.Analyzer
	if cfg.GeminiAPIKey != "" {
		analyzer, err = ai.NewAnalyzer(cfg.GeminiAPIKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize AI analyzer: %v", err)
		}
	}

	return &Bot{
		repo:     repo,
		parser:   feed.NewParser(rssFeedURL),
		notifier: notifier,
		analyzer: analyzer,
	}, nil
}

func (b *Bot) Close() {
	b.repo.Close()
	b.notifier.Close()
	if b.analyzer != nil {
		b.analyzer.Close()
	}
}

func (b *Bot) checkFeed() {
	items, err := b.parser.Parse()
	if err != nil {
		log.Printf("Failed to parse feed: %v", err)
		return
	}

	for i := len(items) - 1; i >= 0; i-- {
		item := items[i]

		exists, err := b.repo.Exists(item.Link)
		if err != nil {
			log.Printf("Failed to check notification: %v", err)
			continue
		}

		if exists {
			continue
		}

		content := item.Description
		if content == "" {
			content = item.Content
		}

		if err := b.notifier.Send(item.Title, item.Link, content, item.PublishedParsed); err != nil {
			log.Printf("Failed to send notification: %v", err)
			continue
		}

		if b.analyzer != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			cleanedContent := cleanContentForAI(content)
			analysis, err := b.analyzer.Analyze(ctx, cleanedContent)
			if err != nil {
				log.Printf("Failed to generate AI analysis: %v", err)
			} else if analysis != "" {
				if err := b.notifier.SendAIAnalysis(analysis); err != nil {
					log.Printf("Failed to send AI analysis: %v", err)
				}
			}
		}

		if err := b.repo.Save(item.Link); err != nil {
			log.Printf("Failed to save notification: %v", err)
			continue
		}

		log.Printf("Notification sent: %s", item.Title)
	}
}

func cleanContentForAI(content string) string {
	text := content

	shortcodePatterns := []string{
		`\[carousel\].*?\[/carousel\]`,
		`\[img\].*?\[/img\]`,
		`\[previewyoutube\].*?\[/previewyoutube\]`,
		`\[youtube\].*?\[/youtube\]`,
		`\[video\].*?\[/video\]`,
		`\[list\].*?\[/list\]`,
		`\[olist\].*?\[/olist\]`,
		`\[code\].*?\[/code\]`,
		`\[quote\].*?\[/quote\]`,
		`\[spoiler\].*?\[/spoiler\]`,
		`\[url[^\]]*\].*?\[/url\]`,
	}

	for _, pattern := range shortcodePatterns {
		re := regexp.MustCompile(`(?is)` + pattern)
		text = re.ReplaceAllString(text, "")
	}

	selfClosingShortcodes := []string{
		`\[img[^\]]*\]`,
		`\[previewyoutube[^\]]*\]`,
		`\[youtube[^\]]*\]`,
		`\[video[^\]]*\]`,
	}

	for _, pattern := range selfClosingShortcodes {
		re := regexp.MustCompile(`(?is)` + pattern)
		text = re.ReplaceAllString(text, "")
	}

	text = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</li>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)<li>`).ReplaceAllString(text, "• ")
	text = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)<[^>]+>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	return strings.TrimSpace(text)
}

func (b *Bot) Start() error {
	if err := b.notifier.Connect(); err != nil {
		return err
	}

	log.Println("Bot started successfully. Press CTRL+C to exit.")

	b.checkFeed()

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Checking for updates...")
		b.checkFeed()
	}

	return nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bot, err := NewBot(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize bot: %v", err)
	}
	defer bot.Close()

	go func() {
		if err := bot.Start(); err != nil {
			log.Fatalf("Bot error: %v", err)
		}
	}()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("Shutting down...")
}
