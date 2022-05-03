package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/db"
)

func EmojiUpdate(session *discordgo.Session, event *discordgo.GuildEmojisUpdate) {
	db.UpdateDiscordEmoji(&event.GuildID, event.Emojis)
}
