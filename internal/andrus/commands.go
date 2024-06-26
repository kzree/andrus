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
	isInChannel := a.checkIfInVoiceChannel(m)
	if isInChannel {
		return
	}

	vs, err := a.findVoiceChannel(m)

	if err != nil || vs == nil {
		a.sendMessage(m.ChannelID, "You must be in a voice channel to use this command!")
		return
	}

	_, err = a.discord.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to join voice channel")
	}
}

func (a *Andrus) leaveCommandHandler(m *discordgo.MessageCreate) {
	vs, err := a.findVoiceChannel(m)

	if err != nil || vs == nil {
		a.sendMessage(m.ChannelID, "You must be in a voice channel to use this command!")
	}

	vc := a.getCurrentVoiceConnection(vs.GuildID)
	if vc == nil {
		a.logger.Error().Msg("failed to find voice connection")
		return
	}

	err = vc.Disconnect()
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to leave voice channel")
	}
}
