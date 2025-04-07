package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	handlers "Menex-bot/m/cmd"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get bot token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Bot token not found. Set BOT_TOKEN as an environment variable.")
	}

	// Initialize bot session
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Error creating Discord session: %v", err)
	}

	// ✅ Enable Message Content Intent
	bot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	bot.AddHandler(handlers.HandleMessage)

	// Open the connection
	err = bot.Open()
	if err != nil {
		log.Fatalf("Error opening connection: %v", err)
	}

	fmt.Println("Bot is running... Press CTRL+C to exit.")

	// Wait for termination signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down bot...")
	bot.Close()
}
