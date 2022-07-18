package handler

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
)

func MessageReactionAdd(session *discordgo.Session, orgReaction *discordgo.MessageReactionAdd) {
	var start time.Time
	if config.CurrentConfig.Debug {
		start = time.Now()
	}

	db.CacheReacted(&orgReaction.MessageID, &orgReaction.Emoji.ID)

	// Ignore all reactions from blacklisted user
	for _, v := range config.CurrentConfig.UserBlacklist {
		if orgReaction.UserID == v {
			return
		}
	}

	if orgReaction.Emoji.ID != "" {
		count, err := db.GetReacted(session, &orgReaction.ChannelID, &orgReaction.MessageID, &orgReaction.Emoji.ID)
		if err != nil {
			log.Print("WARN: failed to get IsReacted: ", err)
			return
		}
		if count <= 0 {
			log.Print("WARN: failed to find added reaction")
			return
		}
		if count == 1 {
			db.AddLog(&orgReaction.GuildID, db.REACTNEW, &orgReaction.Emoji.ID, orgReaction.Member.User.Bot, &orgReaction.UserID, &orgReaction.ChannelID, &orgReaction.MessageID)
		} else {
			db.AddLog(&orgReaction.GuildID, db.REACTADD, &orgReaction.Emoji.ID, orgReaction.Member.User.Bot, &orgReaction.UserID, &orgReaction.ChannelID, &orgReaction.MessageID)
		}
	}

	if config.CurrentConfig.Debug {
		log.Print("Reaction processed in ", time.Since(start).Milliseconds(), "ms.")
	}
}
