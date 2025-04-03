package handlers

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// HandleMenexCommand processes the !menex command
func HandleMenexCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Check if the message content is the command we're looking for
	if m.Content == "!menex" {
		// Send the first message immediately
		s.ChannelMessageSend(m.ChannelID, "!samurugerar")

		// Wait for 5 seconds
		time.Sleep(5 * time.Second)

		// Send the second message after the delay
		s.ChannelMessageSend(m.ChannelID, "esse israel é uma resenha")
	}
}
