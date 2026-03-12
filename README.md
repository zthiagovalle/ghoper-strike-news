# 🐹 Gopher Strike News

<div align="center">

**Your tactical CS2 news analyst powered by Go & AI**

*Not just news — intelligence for e-sports enthusiasts and skin traders*

[![Go Version](https://img.shields.io/badge/Go-1.26.1-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![SQLite](https://img.shields.io/badge/SQLite-Modernc.org-003B57?style=for-the-badge&logo=sqlite)](https://pkg.go.dev/modernc.org/sqlite)
[![Discord](https://img.shields.io/badge/Discord-DiscordGo-5865F2?style=for-the-badge&logo=discord)](https://github.com/bwmarrin/discordgo)
[![Google Gemini](https://img.shields.io/badge/Gemini-3.1_Flash_Lite-4285F4?style=for-the-badge&logo=google)](https://ai.google.dev/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)

</div>

---

## 🎯 What is Gopher Strike News?

Gopher Strike News is a **Discord bot** that monitors Counter-Strike 2 updates from the official Steam RSS feed. But it doesn't just forward news — it acts as your personal **e-sports analyst** and **Steam market advisor**.

When Valve releases patch notes, our tactical Gopher:
- 📰 Delivers clean, formatted updates to your Discord channel
- 🤖 Runs AI-powered analysis using Google Gemini
- 💡 Highlights key changes affecting gameplay meta
- 💰 Identifies market opportunities for skins and cases

**Think of it as having a veteran trader and analyst in your Discord server 24/7.**

---

## ✨ Features

### 🔔 24/7 RSS Monitoring
Continuously watches the official [Steam CS2 News Feed](https://store.steampowered.com/feeds/news/app/730/) every 15 minutes. Never miss an update again.

### 🧹 Smart Content Cleaning
Automatically strips:
- HTML tags and entities
- Steam shortcodes (`[carousel]`, `[img]`, `[previewyoutube]`, etc.)
- Formatting artifacts

Delivering only **clean, readable content** to your channel.

### 🤖 AI-Powered Triple Analysis
Every update triggers a Gemini 3.1 Flash-Lite analysis covering:

| Pillar | What You Get |
|--------|--------------|
| ⭐ **Update Highlight** | The most relevant changes identified |
| 🎯 **Meta Impact** | How it affects gameplay, strategies, and competitive scene |
| 💰 **Market Vision** | Profit opportunities, skin investments, case economics |

### 📊 Dual Embed System
Two beautifully formatted Discord embeds per update:

| Embed | Color | Content |
|-------|-------|---------|
| 📰 **News** | Orange | Original patch notes with proper formatting |
| 🤖 **AI Analysis** | Purple | Gemini-generated insights in Portuguese |

### 💾 Persistent Tracking
SQLite database ensures no duplicate notifications. Even after restarts, the bot remembers what it already sent.

---

## 🛠 Tech Stack

| Technology | Purpose |
|------------|---------|
| ![Go](https://img.shields.io/badge/Go-1.26.1-00ADD8?logo=go) | Core runtime |
| ![DiscordGo](https://img.shields.io/badge/DiscordGo-Latest-5865F2?logo=discord) | Discord API integration |
| ![Gemini](https://img.shields.io/badge/Gemini_SDK-Latest-4285F4?logo=google) | AI analysis engine |
| ![SQLite](https://img.shields.io/badge/SQLite-Pure_Go-003B57?logo=sqlite) | Zero-dependency persistence |
| ![Docker](https://img.shields.io/badge/Docker-Multi_Stage-2496ED?logo=docker) | Containerized deployment |

---

## 🚀 Quick Start

### Prerequisites

- Go 1.26.1+ (for local development)
- Docker (for containerized deployment)
- Discord Bot Token
- Google Gemini API Key

### 1️⃣ Clone & Configure

```bash
git clone https://github.com/zthiagovalle/ghoper-strike-news.git
cd ghoper-strike-news

# Copy environment template
cp .env.example .env
```

### 2️⃣ Set Up Environment Variables

Edit `.env` with your credentials:

```env
# Discord Configuration
DISCORD_TOKEN=your_discord_bot_token_here
CHANNEL_ID=your_target_channel_id_here

# AI Configuration
GEMINI_API_KEY=your_google_gemini_api_key_here

# Database (optional - defaults to cs2bot.db)
DATABASE_URL=/data/cs2bot.db
```

<details>
<summary>📖 How to get these credentials</summary>

**Discord Bot Token:**
1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Create a New Application
3. Navigate to Bot → Add Bot
4. Copy the Token
5. Enable Message Embed Links permission
6. Use OAuth2 → URL Generator to invite bot to your server

**Channel ID:**
1. Enable Developer Mode in Discord (User Settings → Advanced)
2. Right-click your channel → Copy ID

**Gemini API Key:**
1. Visit [Google AI Studio](https://aistudio.google.com/)
2. Create an API key
3. Copy the key

</details>

### 3️⃣ Run Locally

```bash
# Install dependencies
go mod tidy

# Run the bot
go run ./cmd/bot
```

### 4️⃣ Deploy with Docker

```bash
# Build the image
docker build -t gopher-strike-news .

# Run the container
docker run -d \
  --name gopher-strike-news \
  -e DISCORD_TOKEN=your_token \
  -e CHANNEL_ID=your_channel_id \
  -e GEMINI_API_KEY=your_gemini_key \
  -v gopher-data:/data \
  gopher-strike-news
```

### ☁️ Deploy to Railway

[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/your-template)

1. Connect your GitHub repository
2. Set environment variables in Railway dashboard
3. Add a persistent volume at `/data`
4. Deploy!

---

## 📁 Project Structure

```
ghoper-strike-news/
├── cmd/bot/main.go           # Application entry point
├── internal/
│   ├── ai/analyzer.go        # Gemini AI integration
│   ├── config/config.go      # Environment configuration
│   ├── database/database.go  # SQLite persistence
│   ├── discord/discord.go    # Discord embed builder
│   └── feed/feed.go          # RSS parser
├── Dockerfile                # Multi-stage build
├── .env.example              # Configuration template
└── README.md
```

---

## 🎮 Example Output

When CS2 releases an update, your Discord channel receives:

**Embed 1 - News (Orange)**
```
> [RELEASE] CS2 Update - March 2026
> 
> • Fixed collision issues on Mirage
> • Updated weapon pricing for M4A1-S
> • New case: Operation Phoenix 2
```

**Embed 2 - AI Analysis (Purple)**
```
⭐ Destaque da Atualização:
A redução de preço do M4A1-S é a mudança mais impactante...

🎯 Impacto no Jogo:
Espera-se um aumento de 40% na utilização do M4A1-S...

💰 Visão de Mercado:
Oportunidade de compra: skins M4A1-S antes do hype...
```

---

## 🤝 Contributing

Contributions are welcome! Feel free to:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👤 Author

<div align="center">

**Thiago Valle**

[![GitHub](https://img.shields.io/badge/GitHub-zthiagovalle-181717?style=for-the-badge&logo=github)](https://github.com/zthiagovalle)

*Built with 💛 for the CS2 community*

</div>

---

<div align="center">

**⭐ If this project helped you, give it a star! ⭐**

*GG WP 🐹*

</div>
