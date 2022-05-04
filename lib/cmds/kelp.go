package cmds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Kelp = "kelp"

func KelpCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	if config.CurrentConfig.Kelp.Description != "" || config.CurrentConfig.Kelp.Image != "" {
		msg := embed.NewEmbed(session, orgMsg)
		msg.Title = "Kelp"
		if config.CurrentConfig.Kelp.Description != "" {
			msg.Description = config.CurrentConfig.Kelp.Description
		}
		if config.CurrentConfig.Kelp.Image != "" {
			msg.Image = &discordgo.MessageEmbedImage{
				URL: config.CurrentConfig.Kelp.Image,
			}
		}
		ReplyEmbed(session, orgMsg, msg)
	}
}
