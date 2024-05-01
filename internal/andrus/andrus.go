package andrus

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"kzree.com/andrus/internal/logger"
)

type Andrus struct {
	ds     *discordgo.Session
	logger *zerolog.Logger
}

func New(token string, env string) (*Andrus, error) {
	l := logger.New(env)
	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		l.Error().Err(err).Msg("failed to create discord session")
		return nil, err
	}

	l.Info().Msg("created Discord session")
	return &Andrus{ds: ds, logger: l}, nil
}

func (a *Andrus) Run() {
	a.logger.Info().Msg("starting Andrus service")
}
