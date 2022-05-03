package cmds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Ping = "ping"

func PingCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	embedMsg := embed.NewEmbed(session, orgMsg)
	embedMsg.Title = "Pong!"
	ReplyEmbed(session, orgMsg, embedMsg)
}
