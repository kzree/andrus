package andrus

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"kzree.com/andrus/internal/queue"
)

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
	_, err = a.yt.DownloadMedia(a.current)
	a.logger.Info().Msg("finished media download")
	a.logger.Debug().Interface("media", a.current).Msg("playing media")
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to download media")
		return
	}

	a.sendMessage(m.ChannelID, fmt.Sprintf("Playing now! - **%s**", a.current.Title))

	options := dca.StdEncodeOptions
	options.BufferedFrames = 100
	options.FrameDuration = 20
	options.CompressionLevel = 5
	options.Bitrate = 96

	a.logger.Info().Str("filePath", *a.current.FilePath).Msg("creating encoding session")
	encodeSession, err := dca.EncodeFile(*a.current.FilePath, options)
	if err != nil {
		a.logger.Error().Err(err).Msg("failed to create encoding session")
		return
	}
	defer encodeSession.Cleanup()

	time.Sleep(500 * time.Millisecond)

	done := make(chan error)
	dca.NewStream(encodeSession, vc, done)

	select {
	case err := <-done:
		if err != nil && err != io.EOF {
			a.logger.Error().Err(err).Msg("failed to stream audio")
			return
		}
	}
}
