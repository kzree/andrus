package andrus

func (a *Andrus) sendMessage(channelID string, msg string) {
	a.logger.Info().Str("channel", channelID).Str("msg", msg).Msg("sending message")
	a.discord.ChannelMessageSend(channelID, msg)
}
