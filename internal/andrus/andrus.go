package andrus

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
	"kzree.com/andrus/internal/logger"
)

type Andrus struct {
	discord *discordgo.Session
	logger  *zerolog.Logger
}

func New(token string, env string) (*Andrus, error) {
	l := logger.New(env)
	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		l.Error().Err(err).Msg("failed to create discord session")
		return nil, err
	}

	l.Info().Msg("created Discord session")
	return &Andrus{discord: ds, logger: l}, nil
}

func (a *Andrus) Run() {
	a.logger.Info().Msg("starting Andrus service")

	err := a.discord.Open()
	if err != nil {
		a.logger.Fatal().Err(err).Msg("failed to open discord session")
	}

	a.logger.Info().Msg("started Andrus service successfully")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	a.logger.Info().Msg("shutting down Andrus service")
	a.discord.Close()
}
