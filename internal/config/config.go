package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string
	ChannelID    string
	GeminiAPIKey string
	DatabaseURL  string
}

const defaultDatabaseURL = "cs2bot.db"

func Load() (*Config, error) {
	_ = godotenv.Load()

	token := os.Getenv("DISCORD_TOKEN")
	channelID := os.Getenv("CHANNEL_ID")
	geminiKey := os.Getenv("GEMINI_API_KEY")
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		databaseURL = defaultDatabaseURL
	}

	if token == "" || channelID == "" {
		return nil, fmt.Errorf("DISCORD_TOKEN and CHANNEL_ID are required")
	}

	return &Config{
		DiscordToken: token,
		ChannelID:    channelID,
		GeminiAPIKey: geminiKey,
		DatabaseURL:  databaseURL,
	}, nil
}
