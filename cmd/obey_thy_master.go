package cmd

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// messageCreate is called whenever a message is created in any channel the bot has access to
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Handle the !menex command
	if m.Content == "!menex" {
		// Send the first message immediately
		s.ChannelMessageSend(m.ChannelID, "!samurugerar")

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)

		// Send the second message after the delay
		s.ChannelMessageSend(m.ChannelID, "esse israel é uma resenha")
	}

	// You can add more command handlers here
	// For example:
	// if strings.HasPrefix(m.Content, "!help") {
	//     s.ChannelMessageSend(m.ChannelID, "Available commands: !menex")
	// }
}
