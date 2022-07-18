package handler

import (
	"log"
	"runtime/debug"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/cmds"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
)

func MessageCreate(session *discordgo.Session, orgMsg *discordgo.MessageCreate) {
	var start time.Time
	if config.CurrentConfig.Debug {
		start = time.Now()
	}

	guild := db.LoadGuild(&orgMsg.GuildID)
	defer func() {
		if err := recover(); err != nil {
			log.Print("Oops, ", err)
			debug.PrintStack()
		}
	}()

	db.MessageRecieved(orgMsg.Message)

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if orgMsg.Author.ID == session.State.User.ID || orgMsg.Content == "" {
		return
	}

	// Ignore all messages from blacklisted user
	for _, v := range config.CurrentConfig.UserBlacklist {
		if orgMsg.Author.ID == v {
			return
		}
	}

	isCmd := false
	var trimedMsg string
	if strings.HasPrefix(orgMsg.Content, guild.Prefix) {
		isCmd = true
		trimedMsg = strings.TrimPrefix(orgMsg.Content, guild.Prefix)
	} else if strings.HasPrefix(orgMsg.Content, session.State.User.Mention()) {
		isCmd = true
		trimedMsg = strings.TrimPrefix(orgMsg.Content, session.State.User.Mention())
		trimedMsg = strings.TrimPrefix(trimedMsg, " ")
	}
	if isCmd {
		cmds.HandleCmd(session, orgMsg, guild, &trimedMsg)
		if config.CurrentConfig.Debug {
			log.Print("Processed in ", time.Since(start).Milliseconds(), "ms.")
		}
		return
	}
	if len(orgMsg.GetCustomEmojis()) != 0 {
		m := make(map[string]struct{})
		for _, v := range orgMsg.GetCustomEmojis() {
			m[v.ID] = struct{}{}
		}
		for i := range m {
			db.AddLog(&orgMsg.GuildID, db.MSG, &i, orgMsg.Author.Bot, &orgMsg.Author.ID, &orgMsg.ChannelID, &orgMsg.ID)
		}
	}
	if config.CurrentConfig.Debug {
		log.Print("Message processed in ", time.Since(start).Milliseconds(), "ms.")
	}
}
