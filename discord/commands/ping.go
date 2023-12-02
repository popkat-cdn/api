package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Ping(s *discordgo.Session, i *discordgo.ApplicationCommandInteractionData) {
	// Calculate latency
	latency := s.HeartbeatLatency().Milliseconds()
	response := fmt.Sprintf("Pong! Current latency is %d ms", latency)

	// Send the response
	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: response,
		},
	})
}
