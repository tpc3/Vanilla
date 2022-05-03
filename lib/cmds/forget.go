package cmds

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
)

const Forget = "forget"

func ForgetCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🤔")
	emojisID := strings.Split(*message, " ")
	for _, v := range emojisID {
		_, err := db.DeleteEmoji(&orgMsg.GuildID, v)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		_, err = db.DeleteLogEmoji(&orgMsg.GuildID, v)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
