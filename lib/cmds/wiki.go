package cmds

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Wiki = "wiki"

func WikiExportUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = "`" + err.Error() + "`\n"
		msg.Color = embed.ColorPink
	} else {
		msg.Title = config.Lang[guild.Lang].Usage.Title + "wiki export"
	}
	msg.Description += guild.Prefix + Wiki + " export [page] [options]\n" + config.Lang[guild.Lang].Usage.Ranking.Desc
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "[page]",
		Value: config.Lang[guild.Lang].Usage.Ranking.Page,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-i\n--invert",
		Value: config.Lang[guild.Lang].Usage.Ranking.Invert,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-n <number per a page>\n--num <number per a page>",
		Value: config.Lang[guild.Lang].Usage.Ranking.Num,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "-p <period>\n--period <period>",
		Value: config.Lang[guild.Lang].Usage.Ranking.Period,
	})
	ReplyEmbed(session, orgMsg, msg)
}

func WikiCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🤔")
	splitMsg := strings.SplitN(*message, " ", 2)
	switch splitMsg[0] {
	case "import":
		if len(orgMsg.Attachments) != 1 {
			ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.Onefile)
			return
		}
		resp, err := http.Get(orgMsg.Attachments[0].URL)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		if resp.StatusCode != 200 {
			UnknownError(session, orgMsg, &guild.Lang, errors.New("discord return not 200 status: "+resp.Status))
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		rawPage := buf.String()
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🔽")
		emojisRows, err := db.GetEmojis(&orgMsg.GuildID)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
		}
		emojisDB := map[string]emoji{}
		emojisNameDB := map[string]emoji{}
		for emojisRows.Next() {
			e := emoji{}
			emojisRows.Scan(&e.id, &e.name, &e.description)
			emojisDB[e.id] = e
			emojisNameDB[e.name] = e
		}
		matchID := 0
		matchName := 0
		ignored := 0
		update := 0
		splitPage := strings.Split(rawPage, "### ")
		session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🔄")
		for i, v := range splitPage {
			if i == 0 {
				continue
			}
			splitEmoji := strings.SplitN(v, "\n", 3)
			name := splitEmoji[0]
			id := strings.SplitN(splitEmoji[1], "https://cdn.discordapp.com/emojis/", 2)[1]
			id = strings.TrimSuffix(id, ".png)")
			id = strings.TrimSuffix(id, ".gif)")
			var emoji emoji
			var ok bool
			if emoji, ok = emojisDB[id]; ok {
				matchID++
			} else if emoji, ok = emojisNameDB[name]; ok {
				matchName++
			} else {
				ignored++
				continue
			}
			description := strings.TrimSpace(splitEmoji[2])
			if description != emoji.description {
				_, err := db.DeleteEmoji(&orgMsg.GuildID, emoji.id)
				if err != nil {
					UnknownError(session, orgMsg, &guild.Lang, err)
				}
				_, err = db.AddEmoji(&orgMsg.GuildID, emoji.id, emoji.name, description)
				if err != nil {
					UnknownError(session, orgMsg, &guild.Lang, err)
				}
				update++
			}
		}
		msg := embed.NewEmbed(session, orgMsg)
		msg.Title = config.Lang[guild.Lang].Wiki.Title
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  config.Lang[guild.Lang].Wiki.IDMatched,
			Value: strconv.Itoa(matchID),
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  config.Lang[guild.Lang].Wiki.NameMatched,
			Value: strconv.Itoa(matchName),
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  config.Lang[guild.Lang].Wiki.Ignored,
			Value: strconv.Itoa(ignored),
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  config.Lang[guild.Lang].Wiki.Updated,
			Value: strconv.Itoa(update),
		})
		ReplyEmbed(session, orgMsg, msg)
	case "export":
		var (
			invert *struct{}
			period *int64
		)
		if len(splitMsg) == 2 {
			unnamed, err := ParseParam(splitMsg[1], map[string]any{"i": &invert, "p": &period}, map[string]any{"invert": &invert, "period": &period})
			if err != nil {
				RankingUsage(session, orgMsg, guild, err)
				return
			}
			if len(unnamed) != 0 {
				RankingUsage(session, orgMsg, guild, errors.New("too many arguments"))
				return
			}
		}
		if period == nil {
			period = &defaultPeriod
		}
		if *period < 1 {
			RankingUsage(session, orgMsg, guild, errors.New("period must be positive"))
			return
		}
		rows, err := db.GetEmojis(&orgMsg.GuildID)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
		}
		emojisDB := map[string]emoji{}
		for rows.Next() {
			e := emoji{}
			rows.Scan(&e.id, &e.name, &e.description)
			emojisDB[e.id] = e
		}
		rows.Close()
		rows, err = db.GetRanking(&orgMsg.GuildID, 300, 0, *period, invert != nil)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		buff := bytes.Buffer{}
		md := ""
		for rows.Next() {
			var (
				emojiId string
				point   int
			)
			rows.Scan(&emojiId, &point)
			emoji, err := db.GetDiscordEmoji(session, &orgMsg.GuildID, &emojiId)
			if err != nil {
				UnknownError(session, orgMsg, &guild.Lang, err)
				return
			}
			md += "### " + emoji.Name + "\n"
			md += "![" + emoji.Name + "](https://cdn.discordapp.com/emojis/" + emoji.ID
			if emoji.Animated {
				md += ".gif)\n"
			} else {
				md += ".png)\n"
			}
			md += emojisDB[emojiId].description + "\n\n"
		}
		_, err = buff.WriteString(md)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}

		session.ChannelMessageSendComplex(orgMsg.ChannelID, &discordgo.MessageSend{
			Reference: orgMsg.Reference(),
			Files: []*discordgo.File{{
				Name:        "wiki.md",
				ContentType: "text/markdown",
				Reader:      &buff,
			}},
			// AllowedMentions: &discordgo.MessageAllowedMentions{
			// 	RepliedUser: false,
			// },
		})
	default:
		ErrorReply(session, orgMsg, config.Lang[guild.Lang].Error.SubCmd)
		return
	}
}
