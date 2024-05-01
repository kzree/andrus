package andrus

import "github.com/bwmarrin/discordgo"

func (a *Andrus) sendMessage(channelID string, msg string) {
	a.logger.Info().Str("channel", channelID).Str("msg", msg).Msg("sending message")
	_, err := a.discord.ChannelMessageSend(channelID, msg)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to send message")
	}
}

func (a *Andrus) findVoiceChannel(m *discordgo.MessageCreate) (*discordgo.VoiceState, error) {
	c, err := a.discord.State.Channel(m.ChannelID)
	if err != nil {
		a.logger.Error().Err(err).Str("channel", m.ChannelID).Msg("failed to get channel info")
		return nil, err
	}

	g, err := a.discord.State.Guild(c.GuildID)
	if err != nil {
		a.logger.Error().Err(err).Str("guild", c.GuildID).Msg("failed to get guild info")
		return nil, err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			return vs, nil
		}
	}

	return nil, nil
}
