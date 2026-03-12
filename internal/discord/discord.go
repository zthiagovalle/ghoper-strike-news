package discord

import (
	"html"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/microcosm-cc/bluemonday"
)

const (
	embedColor         = 0xFFA500
	aiEmbedColor       = 0x8A2BE2
	cs2Thumbnail       = "https://cdn.akamai.steamstatic.com/steam/apps/730/header.jpg"
	maxContentLength   = 800
	maxAIContentLength = 4000
	truncationSuffix   = "\n\n**... [Click the title to read full notes]**"
	aiTruncationSuffix = "\n\n**[Análise cortada pelo limite do Discord]**"
)

type Notifier struct {
	session   *discordgo.Session
	channelID string
}

func NewNotifier(token, channelID string) (*Notifier, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	session.Identify.Intents = discordgo.IntentsGuildMessages

	return &Notifier{
		session:   session,
		channelID: channelID,
	}, nil
}

func (n *Notifier) Connect() error {
	return n.session.Open()
}

func (n *Notifier) Close() error {
	return n.session.Close()
}

func (n *Notifier) Send(title, link, content string, published *time.Time) error {
	description := buildEmbedDescription(content)
	timestamp := formatTimestamp(published)

	embed := &discordgo.MessageEmbed{
		Title:       title,
		URL:         link,
		Color:       embedColor,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: cs2Thumbnail},
		Description: description,
		Timestamp:   timestamp,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "CS2 News Bot",
		},
	}

	_, err := n.session.ChannelMessageSendEmbed(n.channelID, embed)
	return err
}

func (n *Notifier) SendAIAnalysis(analysis string) error {
	truncated := truncateAIText(analysis)

	embed := &discordgo.MessageEmbed{
		Title:       "🤖 Análise da IA (Gemini)",
		Color:       aiEmbedColor,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: cs2Thumbnail},
		Description: truncated,
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Powered by Google Gemini",
		},
	}

	_, err := n.session.ChannelMessageSendEmbed(n.channelID, embed)
	return err
}

func buildEmbedDescription(content string) string {
	cleaned := cleanHTML(content)
	truncated := truncateText(cleaned)
	return formatBlockQuote(truncated)
}

func cleanHTML(content string) string {
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

	brRegex := regexp.MustCompile(`(?i)<br\s*/?>`)
	text = brRegex.ReplaceAllString(text, "\n")

	liRegex := regexp.MustCompile(`(?i)</li>`)
	text = liRegex.ReplaceAllString(text, "\n")

	text = regexp.MustCompile(`(?i)<li>`).ReplaceAllString(text, "• ")

	pRegex := regexp.MustCompile(`(?i)</p>`)
	text = pRegex.ReplaceAllString(text, "\n")

	text = regexp.MustCompile(`(?i)<p[^>]*>`).ReplaceAllString(text, "")

	text = regexp.MustCompile(`(?i)<ul[^>]*>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)</ul>`).ReplaceAllString(text, "")

	text = regexp.MustCompile(`(?i)<div[^>]*>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(text, "")

	text = regexp.MustCompile(`(?i)<strong[^>]*>`).ReplaceAllString(text, "**")
	text = regexp.MustCompile(`(?i)</strong>`).ReplaceAllString(text, "**")

	text = regexp.MustCompile(`(?i)<b[^>]*>`).ReplaceAllString(text, "**")
	text = regexp.MustCompile(`(?i)</b>`).ReplaceAllString(text, "**")

	text = regexp.MustCompile(`(?i)<i[^>]*>`).ReplaceAllString(text, "*")
	text = regexp.MustCompile(`(?i)</i>`).ReplaceAllString(text, "*")

	text = regexp.MustCompile(`(?i)<em[^>]*>`).ReplaceAllString(text, "*")
	text = regexp.MustCompile(`(?i)</em>`).ReplaceAllString(text, "*")

	text = regexp.MustCompile(`(?i)<a[^>]*href="[^"]*"[^>]*>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)</a>`).ReplaceAllString(text, "")

	policy := bluemonday.StrictPolicy()
	text = policy.Sanitize(text)

	text = html.UnescapeString(text)

	text = regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")
	text = strings.TrimSpace(text)

	return text
}

func truncateText(text string) string {
	if len(text) <= maxContentLength {
		return text
	}

	truncated := text[:maxContentLength]
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxContentLength-100 {
		truncated = truncated[:lastNewline]
	}

	return truncated + truncationSuffix
}

func truncateAIText(text string) string {
	if len(text) <= maxAIContentLength {
		return text
	}

	truncated := text[:maxAIContentLength]
	lastNewline := strings.LastIndex(truncated, "\n")
	if lastNewline > maxAIContentLength-200 {
		truncated = truncated[:lastNewline]
	}

	return truncated + aiTruncationSuffix
}

func formatBlockQuote(text string) string {
	lines := strings.Split(text, "\n")
	quotedLines := make([]string, len(lines))
	for i, line := range lines {
		quotedLines[i] = "> " + line
	}
	return strings.Join(quotedLines, "\n")
}

func formatTimestamp(t *time.Time) string {
	if t == nil {
		return time.Now().Format(time.RFC3339)
	}
	return t.Format(time.RFC3339)
}
