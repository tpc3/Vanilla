package cmds

import (
	"log"
	"runtime/debug"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/embed"
)

func ReplyEmbed(session *discordgo.Session, orgMsg *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	reply := discordgo.MessageSend{}
	reply.Embed = embed
	reply.Reference = orgMsg.Reference()
	reply.AllowedMentions = &discordgo.MessageAllowedMentions{
		RepliedUser: false,
	}
	_, err := session.ChannelMessageSendComplex(orgMsg.ChannelID, &reply)
	if err != nil {
		log.Print("Failed to send reply: ", err)
	}
}

func ErrorReply(session *discordgo.Session, orgMsg *discordgo.MessageCreate, description string) {
	msgEmbed := embed.NewEmbed(session, orgMsg)
	msgEmbed.Title = "Error"
	msgEmbed.Color = embed.ColorPink
	msgEmbed.Description = description
	ReplyEmbed(session, orgMsg, msgEmbed)
}

func UnknownError(session *discordgo.Session, orgMsg *discordgo.MessageCreate, lang *string, err error) {
	debug.PrintStack()
	msgEmbed := embed.NewEmbed(session, orgMsg)
	msgEmbed.Title = config.Lang[*lang].Error.UnknownTitle
	msgEmbed.Description = config.Lang[*lang].Error.UnknownDesc + "\n`" + err.Error() + "`"
	msgEmbed.Color = embed.ColorPink
	ReplyEmbed(session, orgMsg, msgEmbed)
}

func HandleCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	splitMsg := strings.SplitN(*message, " ", 2)
	var param string
	if len(splitMsg) == 2 {
		param = splitMsg[1]
	} else {
		param = ""
	}
	switch splitMsg[0] {
	case Ping:
		PingCmd(session, orgMsg, guild, &param)
	case Ranking:
		RankingCmd(session, orgMsg, guild, &param)
	case Sync:
		SyncCmd(session, orgMsg, guild, &param)
	case Forget:
		ForgetCmd(session, orgMsg, guild, &param)
	case Config:
		ConfigCmd(session, orgMsg, guild, &param)
	case Wiki:
		WikiCmd(session, orgMsg, guild, &param)
	case Desc:
		DescCmd(session, orgMsg, guild, &param)
	case Help:
		HelpCmd(session, orgMsg, guild, &param)
	case Kelp:
		KelpCmd(session, orgMsg, guild, &param)
	case Export:
		ExportCmd(session, orgMsg, guild, &param)
	case Transfer:
		TransferCmd(session, orgMsg, guild, &param)
	default:
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.NoCmd)
	}
}
