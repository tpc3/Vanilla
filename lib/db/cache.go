package db

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
)

var (
	guildCache map[string]*config.Guild
	emojiCache map[string]*map[string]*discordgo.Emoji
)

func init() {
	guildCache = map[string]*config.Guild{}
	emojiCache = map[string]*map[string]*discordgo.Emoji{}
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
