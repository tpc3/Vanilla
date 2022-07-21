package cmds

import (
	"bytes"
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Export = "export"

func ExportUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = "`" + err.Error() + "`\n"
		msg.Color = embed.ColorPink
	} else {
		msg.Title = config.Lang[guild.Lang].Usage.Title + "wiki export"
	}
	msg.Description += guild.Prefix + Export + " [options]\n" + config.Lang[guild.Lang].Usage.Ranking.Desc
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-i\n--invert",
		Value: config.Lang[guild.Lang].Usage.Ranking.Invert,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-p <period>\n--period <period>",
		Value: config.Lang[guild.Lang].Usage.Ranking.Period,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-b\n--bots",
		Value: config.Lang[guild.Lang].Usage.Ranking.Bots,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-o\n--only-bots",
		Value: config.Lang[guild.Lang].Usage.Ranking.OnlyBots,
	})
	ReplyEmbed(session, orgMsg, msg)
}

func ExportCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🤔")
	splitMsg := strings.SplitN(*message, " ", 2)
	var (
		invert   *struct{}
		period   *int64
		bots     *struct{}
		onlyBots *struct{}
	)
	if len(splitMsg) == 2 {
		unnamed, err := ParseParam(splitMsg[1], map[string]any{"i": &invert, "p": &period, "b": &bots, "o": &onlyBots},
			map[string]any{"invert": &invert, "period": &period, "bots": &bots, "only-bots": &onlyBots})
		if err != nil {
			ErrorReply(session, orgMsg, err.Error())
			return
		}
		if len(unnamed) != 0 {
			ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.Syntax)
			return
		}
	}
	if period == nil {
		period = &defaultPeriod
	}
	if *period < 1 {
		WikiExportUsage(session, orgMsg, guild, errors.New("period must be positive"))
		return
	}
	rows, err := db.GetRanking(&orgMsg.GuildID, 300, 0, *period, invert != nil, (onlyBots == nil), (bots != nil || onlyBots != nil))
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	result := ""
	for rows.Next() {
		var (
			emojiId     string
			emojiName   string
			description string
			point       int
		)
		rows.Scan(&emojiId, &emojiName, &description, &point)
		emoji, err := db.GetDiscordEmoji(session, &orgMsg.GuildID, &emojiId)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		if emoji == nil {
			ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.DeletedEmojiFound)
			return
		}
		result += emoji.Name + ": " + "https://cdn.discordapp.com/emojis/" + emoji.ID
		if emoji.Animated {
			result += ".gif\n"
		} else {
			result += ".webp\n"
		}
	}
	buff := bytes.Buffer{}
	_, err = buff.WriteString(result)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}

	session.ChannelMessageSendComplex(orgMsg.ChannelID, &discordgo.MessageSend{
		Reference: orgMsg.Reference(),
		Files: []*discordgo.File{{
			Name:        "emojis.txt",
			ContentType: "text/plain",
			Reader:      &buff,
		}},
		// AllowedMentions: &discordgo.MessageAllowedMentions{
		// 	RepliedUser: false,
		// },
	})
}
