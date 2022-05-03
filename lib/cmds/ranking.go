package cmds

import (
	"errors"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Ranking = "ranking"

var (
	defaultNum    int
	defaultPeriod int64
)

func init() {
	defaultNum = 5
	defaultPeriod = 2592000
}

func RankingUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = err.Error() + "\n"
		msg.Color = embed.ColorPink
	} else {
		msg.Title = config.Lang[guild.Lang].Usage.Title + "ranking"
	}
	msg.Description += "`" + guild.Prefix + Ranking + " [page] [options]`\n" + config.Lang[guild.Lang].Usage.Ranking.Desc
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

func RankingCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	var (
		page   int
		invert *struct{}
		num    *int
		period *int64
	)
	unnamed, err := ParseParam(*message, map[string]any{"i": &invert, "n": &num, "p": &period}, map[string]any{"invert": &invert, "num": &num, "period": &period})
	if err != nil {
		RankingUsage(session, orgMsg, guild, err)
		return
	}
	if len(unnamed) == 0 {
		page = 1
	} else if len(unnamed) == 1 {
		page, err = strconv.Atoi(unnamed[0])
		if err != nil {
			RankingUsage(session, orgMsg, guild, errors.New("page must be number"))
			return
		}
	} else {
		RankingUsage(session, orgMsg, guild, errors.New("too many arguments"))
		return
	}
	if num == nil {
		num = &defaultNum
	}
	if period == nil {
		period = &defaultPeriod
	}
	if page < 1 {
		RankingUsage(session, orgMsg, guild, errors.New("page must be positive"))
		return
	}
	if *num < 1 {
		RankingUsage(session, orgMsg, guild, errors.New("num must be positive"))
		return
	}
	if *num > 25 {
		RankingUsage(session, orgMsg, guild, errors.New("num must be 25 or less"))
		return
	}
	if *period < 1 {
		RankingUsage(session, orgMsg, guild, errors.New("period must be positive"))
		return
	}
	rank := *num * (page - 1)
	rows, err := db.GetRanking(&orgMsg.GuildID, *num, rank, *period, invert != nil)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	msg := embed.NewEmbed(session, orgMsg)
	msg.Title = config.Lang[guild.Lang].Ranking
	for rows.Next() {
		rank++
		var (
			emojiId string
			point   int
		)
		rows.Scan(&emojiId, &point)
		field := discordgo.MessageEmbedField{}
		emoji, err := db.GetDiscordEmoji(session, &orgMsg.GuildID, &emojiId)
		if err != nil {
			UnknownError(session, orgMsg, &guild.Lang, err)
			return
		}
		field.Name = strconv.Itoa(rank) + ". " + emoji.Name
		field.Value = emoji.MessageFormat() + " " + strconv.FormatInt(int64(point), 10) + "pt"
		msg.Fields = append(msg.Fields, &field)
	}
	ReplyEmbed(session, orgMsg, msg)
}
