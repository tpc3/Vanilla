package handler

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
)

func MessageReactionAdd(session *discordgo.Session, orgReaction *discordgo.MessageReactionAdd) {
	if config.CurrentConfig.Debug {
		start := time.Now()
		defer func() {
			log.Print("Reaction processed in ", time.Since(start).Milliseconds(), "ms.")
		}()
	}

	// Ignore all reactions from blacklisted user
	for _, v := range config.CurrentConfig.UserBlacklist {
		if orgReaction.UserID == v {
			return
		}
	}

	if orgReaction.Emoji.ID != "" {
		msg, err := session.State.Message(orgReaction.ChannelID, orgReaction.MessageID)
		if err != nil {
			msg, err = session.ChannelMessage(orgReaction.ChannelID, orgReaction.MessageID)
			if err != nil {
				log.Print("WARN: failed to get Message to check reacted: ", err)
				return
			}
			err = session.State.MessageAdd(msg)
			if err != nil {
				log.Print("WARN: failed to state cache Message to check reacted: ", err)
			}
		}
		reactType := db.REACTNEW
		for _, v := range msg.Reactions {
			if v.Emoji.ID == orgReaction.Emoji.ID {
				if v.Count > 1 {
					reactType = db.REACTADD
				}
				break
			}
		}
		db.AddLog(&orgReaction.GuildID, reactType, &orgReaction.Emoji.ID, orgReaction.Member.User.Bot, &orgReaction.UserID, &orgReaction.ChannelID, &orgReaction.MessageID)
	}

}
