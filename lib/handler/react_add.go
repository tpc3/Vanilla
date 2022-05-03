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

	if orgReaction.Emoji.ID != "" {
		msg, err := session.ChannelMessage(orgReaction.ChannelID, orgReaction.MessageID)
		if err != nil {
			log.Print("WARN: failed to get ChannelMessage to get reaction: ", err)
			return
		}

		count := 0
		for _, v := range msg.Reactions {
			if v.Emoji.ID == orgReaction.Emoji.ID {
				count = v.Count
				break
			}
		}
		if count == 0 {
			log.Print("WARN: failed to find added reaction")
			return
		}
		if count == 1 {
			db.AddLog(&orgReaction.GuildID, db.REACTNEW, &orgReaction.Emoji.ID)
		} else {
			db.AddLog(&orgReaction.GuildID, db.REACTADD, &orgReaction.Emoji.ID)
		}
	}

	if config.CurrentConfig.Debug {
		log.Print("Processed in ", time.Since(start).Milliseconds(), "ms.")
	}
}
