package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Menex-bot/m/cmd" // Make sure this import matches your module name and directory structure

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Load bot token from environment variable or replace with your token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Bot token not found. Set BOT_TOKEN as an environment variable.")
	}

	// Create a new Discord session
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session:", err)
	}

	// Register the message handler
	dg.AddHandler(cmd.MessageCreate)

	// Open the websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Fatal("Error opening connection to Discord:", err)
	}

	fmt.Println("Menex-Bot is now running. Press CTRL+C to exit.")

	// Wait here until CTRL+C or another termination signal is received
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	fmt.Println("Shutting down bot...")
	dg.Close()
}
