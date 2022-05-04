package db

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
)

var (
	guildCache map[string]*config.Guild
	emojiCache map[string]*map[string]*discordgo.Emoji
	reactCache map[string]*map[string]int
)

func init() {
	guildCache = map[string]*config.Guild{}
	reactCache = map[string]*map[string]int{}
	emojiCache = map[string]*map[string]*discordgo.Emoji{}
}

func MessageRecieved(message *discordgo.Message) {
	reactCache[message.ID] = &map[string]int{}
	for _, v := range message.Reactions {
		if v.Emoji.ID != "" {
			(*reactCache[message.ID])[v.Emoji.ID] = v.Count
		}
	}
}

func GetReacted(session *discordgo.Session, channelId *string, messageId *string, emojiId *string) (int, error) {
	if reactCache[*messageId] == nil {
		reactCache[*messageId] = &map[string]int{}
		msg, err := session.ChannelMessage(*channelId, *messageId)
		if err != nil {
			return -1, err
		}
		MessageRecieved(msg)
	}
	return (*reactCache[*messageId])[*emojiId], nil
}
func CacheReacted(messageId *string, emojiId *string) {
	if reactCache[*messageId] != nil {
		(*reactCache[*messageId])[*emojiId]++
	}
}
func GetDiscordEmojis(session *discordgo.Session, guildId *string) (*map[string]*discordgo.Emoji, error) {
	if emojiCache[*guildId] == nil {
		emojis, err := session.GuildEmojis(*guildId)
		if err != nil {
			return nil, err
		}
		UpdateDiscordEmoji(guildId, emojis)
	}
	return emojiCache[*guildId], nil
}
func GetDiscordEmoji(session *discordgo.Session, guildId *string, emojiId *string) (*discordgo.Emoji, error) {
	emojis, err := GetDiscordEmojis(session, guildId)
	if err != nil {
		return nil, err
	}
	return (*emojis)[*emojiId], nil
}
func UpdateDiscordEmoji(guildId *string, emojis []*discordgo.Emoji) {
	newEmojis := map[string]*discordgo.Emoji{}
	for _, v := range emojis {
		newEmojis[v.ID] = v
	}
	emojiCache[*guildId] = &newEmojis
}
