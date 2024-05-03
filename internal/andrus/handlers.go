package andrus

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (a *Andrus) readyHandler(s *discordgo.Session, _ *discordgo.Ready) {
	a.logger.Info().Msg("updating game status")
	s.UpdateGameStatus(0, "Listening to !play")
}

func (a *Andrus) voiceStateHandler(_ *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.ChannelID == "" {
		a.logger.Info().Msg("leaving voice channel")
	} else {
		a.logger.Info().Str("channel", v.ChannelID).Msg("joining voice channel")
	}
}

func (a *Andrus) createMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!") {
		a.logger.
			Info().
			Str("command", m.Content).
			Interface("author", map[string]any{"id": m.Author.ID, "username": m.Author.Username}).
			Msg("received command")

		command := strings.Split(m.Content, " ")[0]

		switch command {
		case CommandHello:
			a.helloCommandHandler(m)
		case CommandJoin:
			a.joinCommandHandler(m)
		case CommandLeave:
			a.leaveCommandHandler(m)
		case CommandPlay:
			a.playCommandHandler(m)
		default:
			a.logger.Warn().Str("command", m.Content).Msg("failed to find matching command")
		}
	}
}
