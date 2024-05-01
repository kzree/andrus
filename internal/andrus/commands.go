package andrus

import (
	"github.com/bwmarrin/discordgo"
)

const (
	CommandHello = "!hello"
	CommandJoin  = "!join"
	CommandPlay  = "!play"
	CommandStop  = "!stop"
	CommandLeave = "!leave"
	CommandQueue = "!queue"
)

func (a *Andrus) helloCommandHandler(m *discordgo.MessageCreate) {
	a.sendMessage(m.ChannelID, "world!")
}

func (a *Andrus) joinCommandHandler(m *discordgo.MessageCreate) {
	c, err := a.discord.State.Channel(m.ChannelID)
	if err != nil {
		a.logger.Error().Err(err).Str("channel", m.ChannelID).Msg("failed to get channel info")
		return
	}

	g, err := a.discord.State.Guild(c.GuildID)
	if err != nil {
		a.logger.Error().Err(err).Str("guild", c.GuildID).Msg("failed to get guild info")
		return
	}

	foundChannel := false
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			foundChannel = true
			a.logger.Info().Str("channel", vs.ChannelID).Msg("attempting to join voice channel")
			_, err := a.discord.ChannelVoiceJoin(g.ID, vs.ChannelID, false, true)
			if err != nil {
				a.logger.Error().Err(err).Str("channel", c.ID).Msg("failed to join voice channel")
			}

			return
		}
	}

	if !foundChannel {
		a.sendMessage(m.ChannelID, "You must be in a voice channel to use this command!")
	}
}
