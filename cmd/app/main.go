package main

import (
	"flag"

	"kzree.com/andrus/internal/andrus"
)

type Config struct {
	DiscordToken string
	Env          string
}

func main() {
	var cfg Config
	flag.StringVar(&cfg.DiscordToken, "discord-token", "", "Discord authentication token")
	flag.StringVar(&cfg.Env, "env", "dev", "Environment to run the application in")
	flag.Parse()

	app, err := andrus.New(cfg.DiscordToken, cfg.Env)
	if err != nil {
		panic(err)
	}

	app.Run()
}
