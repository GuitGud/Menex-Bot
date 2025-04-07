package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	handlers "Menex-bot/m/cmd"
	"Menex-bot/m/utils"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get bot token and shutdown channel
	token := os.Getenv("BOT_TOKEN")
	notificationChannelID := os.Getenv("SHUTDOWN_CHANNEL_ID") // Add this to your .env
	if token == "" || notificationChannelID == "" {
		log.Fatal("Missing BOT_TOKEN or SHUTDOWN_CHANNEL_ID in environment variables.")
	}

	// Load events at startup
	utils.LoadEvents()

	// Initialize bot session
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// ✅ Enable Message Content Intent
	bot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	// Register message handler
	bot.AddHandler(handlers.HandleMessage)

	// Open connection
	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}
	defer bot.Close()

	// ✅ Send startup message to notification channel (if set)
	if notificationChannelID != "" {
		_, err = bot.ChannelMessageSend(notificationChannelID, "✅ *FENT REACTOR AT 100%* Menex-Bot online.")
		if err != nil {
			log.Println("Failed to send startup message:", err)
		}
	}

	fmt.Println("Menex-Bot is running... Press CTRL+C to exit.")

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	// Send shutdown message to notification channel
	if notificationChannelID != "" {
		_, err = bot.ChannelMessageSend(notificationChannelID, "⚠️ Menex-Bot está desligando...\n*FENT REACTOR POWERING DOWN...* Adeus, humanos.")
		if err != nil {
			log.Println("Failed to send shutdown message:", err)
		}
	}

	// Save any pending data
	handlers.SaveEvents()

	fmt.Println("Shutting down MenexlinkAI... Powering down fent reactor...")
}
