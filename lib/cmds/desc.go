package cmds

import (
	"errors"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Desc = "desc"

func DescCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	splitLine := strings.SplitN(*message, "\n", 2)
	splitMsg := strings.SplitN(splitLine[0], " ", 2)
	if len(splitMsg) != 2 {
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.Syntax)
		return
	}
	id := strings.TrimSpace(splitMsg[1])
	splitName := strings.SplitN(id, ":", 3)
	if len(splitName) == 3 {
		id = strings.TrimSuffix(splitName[2], ">")
	}
	emoji := emoji{}
	_, err := strconv.Atoi(id)
	if err == nil {
		row := db.GetEmoji(&orgMsg.GuildID, id)
		err = row.Scan(&emoji.id, &emoji.name, &emoji.description)
	}
	if err != nil {
		row := db.GetEmojiByName(&orgMsg.GuildID, strings.TrimSpace(splitMsg[1]))
		err = row.Scan(&emoji.id, &emoji.name, &emoji.description)
	}
	if err != nil {
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.EmojiNotFound)
		return
	}
	switch splitMsg[0] {
	case "set":
		emoji.description = splitLine[1]
		_, err := db.DeleteEmoji(&orgMsg.GuildID, emoji.id)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		_, err = db.AddEmoji(&orgMsg.GuildID, emoji.id, emoji.name, emoji.description)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
	case "get":
		msg := embed.NewEmbed(session, orgMsg)
		d, err := session.State.Emoji(orgMsg.GuildID, emoji.id)
		if err != nil {
			if errors.Is(err, discordgo.ErrStateNotFound) {
				ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.DeletedEmojiFound)
			} else {
				UnknownError(session, orgMsg, &guild.Lang, err)
			}
			return
		}
		msg.Title = d.MessageFormat() + " " + emoji.name
		msg.Description = emoji.description
		ReplyEmbed(session, orgMsg, msg)
	default:
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.SubCmd)
		return
	}
}
