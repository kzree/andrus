package andrus

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"kzree.com/andrus/internal/queue"
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

func (a *Andrus) playCommandHandler(m *discordgo.MessageCreate) {
	vs, err := a.findVoiceChannel(m)

	if err != nil || vs == nil {
		a.sendMessage(m.ChannelID, "You must be in a voice channel to use this command!")
		return
	}

	split := strings.Split(m.Content, " ")
	if len(split) < 2 {
		a.sendMessage(m.ChannelID, "You must provide a URL to play!")
		return
	}

	url := split[1]
	// validate that url is valid youtube url
	if !strings.Contains(url, "youtube.com") {
		a.sendMessage(m.ChannelID, "Invalid URL! currently only youtube is supported!")
		return
	}

	vc := a.getCurrentVoiceConnection(vs.GuildID)
	if vc == nil {
		_, err = a.discord.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to join voice channel")
		}
	}

	media := queue.NewMedia(url, "", "", &queue.Requester{ID: m.Author.ID, Username: m.Author.Username})
	if a.current != nil {
		err = a.queue.Add(media)
		if err != nil {
			a.sendMessage(m.ChannelID, "Queue is full!")
			return
		}
		a.sendMessage(m.ChannelID, "Added to queue!")
		return
	}

	a.logger.Info().Interface("media", media).Msg("preparing media for playing")
	a.current = media

	a.logger.Info().Msg("downloading media")
	_, err = a.yt.DownloadMedia(*a.current)
	a.logger.Info().Msg("finished media download")
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to download media")
		return
	}

	a.sendMessage(m.ChannelID, "Playing now!")
}
