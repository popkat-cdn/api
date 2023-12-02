package discord

import (
	"Popkat/discord/commands"
	"Popkat/state"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Setup() {
	state.Discord.Open() // Open Discord Client

	state.Discord.UpdateWatchStatus(0, "Popkat CDN") // Update Presence

	fmt.Println("Discord Bot is now running. Press Ctrl+C to exit.")
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "ping" {
		commands.Ping(s, i.ApplicationCommandData())
	}
}
