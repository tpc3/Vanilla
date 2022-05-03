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

const Config = "config"

func ConfigUsage(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, err error) {
	msg := embed.NewEmbed(session, orgMsg)
	if err != nil {
		msg.Title = config.Lang[guild.Lang].Error.Syntax
		msg.Description = err.Error() + "\n"
		msg.Color = embed.ColorPink
	} else {
		msg.Title = config.Lang[guild.Lang].Usage.Title + "ranking"
	}
	msg.Description += "`" + guild.Prefix + Config + " [<item> <value>]`\n" + config.Lang[guild.Lang].Usage.Config.Desc
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "prefix <prefix>",
		Value: config.Lang[guild.Lang].Usage.Config.Prefix,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "lang <language>",
		Value: config.Lang[guild.Lang].Usage.Config.Lang,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "bots <record bots>",
		Value: config.Lang[guild.Lang].Usage.Config.Bots,
	})
	msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
		Name:  "weight <message> <new reaction> <add reaction>",
		Value: config.Lang[guild.Lang].Usage.Config.Weight,
	})
	ReplyEmbed(session, orgMsg, msg)
}

func ConfigCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	split := strings.SplitN(*message, " ", 2)
	if *message == "" {
		msg := embed.NewEmbed(session, orgMsg)
		msg.Title = config.Lang[guild.Lang].CurrConf
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "prefix",
			Value: guild.Prefix,
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "lang",
			Value: guild.Lang,
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "bots",
			Value: strconv.FormatBool(guild.Recordbots),
		})
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "weight",
			Value: strconv.Itoa(guild.Weight.Message) + " " + strconv.Itoa(guild.Weight.Reactnew) + " " + strconv.Itoa(guild.Weight.Reactadd),
		})
		ReplyEmbed(session, orgMsg, msg)
		return
	}
	if len(split) != 2 {
		ConfigUsage(session, orgMsg, guild, errors.New("not enough arguments"))
		return
	}
	ok := false
	switch split[0] {
	case "prefix":
		guild.Prefix = split[1]
	case "lang":
		_, ok = config.Lang[split[1]]
		if ok {
			guild.Lang = split[1]
		} else {
			ErrorReply(session, orgMsg, "unsupported language")
			return
		}
	case "bots":
		conv, err := strconv.ParseBool(split[1])
		if err != nil {
			ConfigUsage(session, orgMsg, guild, errors.New("failed to parse value"))
			return
		}
		guild.Recordbots = conv
	case "weight":
		weights := strings.Split(split[1], " ")
		if len(weights) != 3 {
			ConfigUsage(session, orgMsg, guild, errors.New("invalid weight arguments length"))
			return
		}
		msg, err := strconv.Atoi(weights[0])
		if err != nil {
			ConfigUsage(session, orgMsg, guild, errors.New("weight must be integer"))
			return
		}
		reactnew, err := strconv.Atoi(weights[1])
		if err != nil {
			ConfigUsage(session, orgMsg, guild, errors.New("weight must be integer"))
			return
		}
		reactadd, err := strconv.Atoi(weights[2])
		if err != nil {
			ConfigUsage(session, orgMsg, guild, errors.New("weight must be integer"))
			return
		}
		guild.Weight.Message = msg
		guild.Weight.Reactnew = reactnew
		guild.Weight.Reactadd = reactadd
	default:
		ConfigUsage(session, orgMsg, guild, errors.New("item not found"))
		return
	}
	err := db.SaveGuild(&orgMsg.GuildID, guild)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "👍")
}
