package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Debug   bool
	Help    string
	Discord struct {
		Token  string
		Status string
	}
	Db struct {
		Kind        string
		Parameter   string
		Tableprefix string
	}
	Guild Guild
}

type Guild struct {
	Prefix     string
	Lang       string
	Recordbots bool
	Weight     struct {
		Message  int
		Reactnew int
		Reactadd int
	}
}

const configFile = "./config.yaml"

var CurrentConfig Config

func init() {
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Config load failed:", err)
	}
	err = yaml.Unmarshal(file, &CurrentConfig)
	if err != nil {
		log.Fatal("Config parse failed:", err)
	}

	//verify
	if CurrentConfig.Debug {
		log.Print("Debug is enabled")
	}
	if CurrentConfig.Discord.Token == "" {
		log.Fatal("Token is empty")
	}
	if CurrentConfig.Db.Tableprefix == "" {
		log.Fatal("Tableprefix is empty")
	}

	loadLang()
}
