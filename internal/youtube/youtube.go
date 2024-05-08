package youtube

import (
	"errors"
	"io"
	"os"
	"regexp"

	kkdai_youtube "github.com/kkdai/youtube/v2"
	"github.com/rs/zerolog"
	"kzree.com/andrus/internal/queue"
)

type Youtube struct {
	logger *zerolog.Logger
	client *kkdai_youtube.Client
}

func New(l *zerolog.Logger) *Youtube {
	logger := l.With().Str("module", "youtube").Logger()
	return &Youtube{
		logger: &logger,
		client: &kkdai_youtube.Client{},
	}
}

func (y *Youtube) extractVideoID(url string) (string, error) {
	regex := regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	matches := regex.FindStringSubmatch(url)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", errors.New("failed to extract video ID")
}

func (y *Youtube) DownloadMedia(m queue.Media) (*string, error) {
	y.logger.Debug().Str("url", m.URL).Msg("starting media download")
	videoID, err := y.extractVideoID(m.URL)
	if err != nil {
		return nil, err
	}
	y.logger.Debug().Str("videoID", videoID).Msg("extracted video ID")

	video, err := y.client.GetVideo(videoID)
	if err != nil {
		return nil, err
	}

	m.Title = video.Title
	y.logger.Debug().Str("title", video.Title).Msg("found video metadata")

	y.logger.Debug().Msg("getting video stream")
	formats := video.Formats.Itag(140) // youtube Itag for m4a audio
	stream, _, err := y.client.GetStream(video, &formats[0])
	if err != nil {
		return nil, err
	}
	defer stream.Close()

	fileName := ".cache/" + video.ID + ".mp3"
	y.logger.Debug().Str("file", fileName).Msg("saving media file to cache")
	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, stream)
	if err != nil {
		return nil, err
	}

	return &fileName, nil
}
