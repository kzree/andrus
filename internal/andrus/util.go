package andrus

func (a *Andrus) sendMessage(channelID string, msg string) {
	a.logger.Info().Str("channel", channelID).Str("msg", msg).Msg("sending message")
	_, err := a.discord.ChannelMessageSend(channelID, msg)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to send message")
	}
}
