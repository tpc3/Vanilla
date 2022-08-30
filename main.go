package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/tpc3/Vanilla/lib/config"
	"github.com/tpc3/Vanilla/lib/db"
	"github.com/tpc3/Vanilla/lib/handler"
)

var discord *discordgo.Session

func main() {
	var err error
	discord, err = discordgo.New("Bot " + config.CurrentConfig.Discord.Token)
	if err != nil {
		log.Fatal("main.go discordgo.New() failed:", err)
	}
	discord.AddHandler(handler.MessageCreate)
	discord.AddHandler(handler.MessageReactionAdd)
	discord.AddHandler(handler.EmojiUpdate)
	err = discord.Open()
	if err != nil {
		log.Print("Discordgo connection failure:", err)
		return
	}
	discord.UpdateGameStatus(0, config.CurrentConfig.Discord.Status)
	log.Print("Started Discordgo")
	defer discord.Close()
	defer db.Close()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Print("Vanilla is gracefully shutdowning!")
}
