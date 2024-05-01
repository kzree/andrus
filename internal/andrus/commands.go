package andrus

import "github.com/bwmarrin/discordgo"

const (
	CommandHello = "!hello"
)

func (a *Andrus) helloCommandHandler(m *discordgo.MessageCreate) {
	a.sendMessage(m.ChannelID, "world!")
}
