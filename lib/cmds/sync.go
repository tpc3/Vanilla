package cmds

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/embed"
)

const Sync = "sync"

type emoji struct {
	id          string
	name        string
	description string
	discord     *discordgo.Emoji
}

func SyncCmd(session *discordgo.Session, orgMsg *discordgo.MessageCreate, guild *config.Guild, message *string) {
	session.MessageReactionAdd(orgMsg.ChannelID, orgMsg.ID, "🤔")
	result := embed.NewEmbed(session, orgMsg)
	result.Title = config.Lang[guild.Lang].Sync.Title
	start := time.Now()
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
	emojisDiscord, err := db.GetDiscordEmojis(session, &orgMsg.GuildID)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
	}
	var addedEmoji []emoji
	var updatedEmoji []emoji
	var movedEmoji []emoji
	var deletedEmoji []emoji
	for _, v := range *emojisDiscord {
		if e, ok := emojisDB[v.ID]; ok {
			delete(emojisNameDB, e.name)
			e.discord = v
			if e.name != v.Name {
				_, err := db.DeleteEmoji(&orgMsg.GuildID, e.id)
				if err != nil {
					UnknownError(session, orgMsg, &guild.Lang, err)
				}
				_, err = db.AddEmoji(&orgMsg.GuildID, e.id, v.Name, e.description)
				if err != nil {
					UnknownError(session, orgMsg, &guild.Lang, err)
				}
				updatedEmoji = append(updatedEmoji, e)
			}
		} else if e, ok := emojisNameDB[v.Name]; ok {
			delete(emojisNameDB, e.name)
			e.discord = v
			_, err := db.DeleteEmoji(&orgMsg.GuildID, e.id)
			if err != nil {
				UnknownError(session, orgMsg, &guild.Lang, err)
			}
			_, err = db.AddEmoji(&orgMsg.GuildID, v.ID, e.name, e.description)
			if err != nil {
				UnknownError(session, orgMsg, &guild.Lang, err)
			}
			movedEmoji = append(movedEmoji, e)
		} else {
			_, err := db.AddEmoji(&orgMsg.GuildID, v.ID, v.Name, "")
			if err != nil {
				UnknownError(session, orgMsg, &guild.Lang, err)
			}
			addedEmoji = append(addedEmoji, emoji{id: v.ID, name: v.Name, description: "", discord: v})
		}
	}
	for _, v := range emojisNameDB {
		deletedEmoji = append(deletedEmoji, v)
	}
	if len(addedEmoji) != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.NewEmoji
		for i, v := range addedEmoji {
			field.Value += v.discord.MessageFormat() + " " + v.name + "\n"
			if len(field.Value) > 900 {
				field.Value += config.Lang[guild.Lang].Sync.OverEmoji + strconv.Itoa(len(addedEmoji)-(i+1))
				break
			}
		}
		result.Fields = append(result.Fields, &field)
	}
	if len(movedEmoji) != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.IDChangedEmoji
		for i, v := range movedEmoji {
			field.Value += v.discord.MessageFormat() + " " + v.name + " " + v.id + " -> " + v.discord.ID + "\n"
			db.ChangeLogID(&orgMsg.GuildID, &v.id, &v.discord.ID)
			if len(field.Value) > 900 {
				field.Value += config.Lang[guild.Lang].Sync.OverEmoji + strconv.Itoa(len(movedEmoji)-(i+1))
				break
			}
		}
		result.Fields = append(result.Fields, &field)
	}
	if len(updatedEmoji) != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.NameChangedEmoji
		for i, v := range updatedEmoji {
			field.Value += v.discord.MessageFormat() + " " + v.name + " -> " + v.discord.Name + "\n"
			if len(field.Value) > 900 {
				field.Value += config.Lang[guild.Lang].Sync.OverEmoji + strconv.Itoa(len(updatedEmoji)-(i+1))
				break
			}
		}
		result.Fields = append(result.Fields, &field)
	}
	if len(deletedEmoji) != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.DeletedEmoji
		cleanCommand := guild.Prefix + Forget
		for i, v := range deletedEmoji {
			field.Value = v.id + " " + v.name + "\n"
			cleanCommand += " " + v.id
			if len(field.Value) > 900 {
				field.Value += config.Lang[guild.Lang].Sync.OverEmoji + strconv.Itoa(len(deletedEmoji)-(i+1))
				break
			}
		}
		field.Value += config.Lang[guild.Lang].Sync.ToCleanEmoji + "`" + cleanCommand + "`"
		result.Fields = append(result.Fields, &field)
	}
	result.Description += "check emoji: " + strconv.FormatInt(time.Since(start).Milliseconds(), 10) + "ms\n"
	start = time.Now()
	var validEmojiID []string
	for i := range *emojisDiscord {
		validEmojiID = append(validEmojiID, i)
	}
	for _, v := range deletedEmoji {
		validEmojiID = append(validEmojiID, v.id)
	}
	updated, err := db.CleanLogEmoji(&orgMsg.GuildID, validEmojiID)
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	if *updated != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.CleanLogTitle
		field.Value = strconv.FormatInt(*updated, 10) + config.Lang[guild.Lang].Sync.CleanLogDesc
		result.Fields = append(result.Fields, &field)
	}
	result.Description += "clean log: " + strconv.FormatInt(time.Since(start).Milliseconds(), 10) + "ms\n"
	start = time.Now()
	updated, err = db.UpdateValue(&orgMsg.GuildID, map[int]int{db.MSG: guild.Weight.Message, db.REACTNEW: guild.Weight.Reactnew, db.REACTADD: guild.Weight.Reactadd})
	if err != nil {
		UnknownError(session, orgMsg, &guild.Lang, err)
		return
	}
	if *updated != 0 {
		field := discordgo.MessageEmbedField{}
		field.Name = config.Lang[guild.Lang].Sync.WeightTitle
		field.Value = strconv.FormatInt(*updated, 10) + config.Lang[guild.Lang].Sync.WeightDesc
		result.Fields = append(result.Fields, &field)
	}
	result.Description += "update weight: " + strconv.FormatInt(time.Since(start).Milliseconds(), 10) + "ms\n"
	ReplyEmbed(session, orgMsg, result)
}
