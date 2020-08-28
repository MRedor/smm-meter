package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type сonfig struct {
	Google   *oauth2.Config
	DB DBConfig `yaml:"database"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

var Cfg сonfig

func initGoogle() {
	b, err := ioutil.ReadFile("src/config/second_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	Cfg.Google, err = google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
}

func initDB() {
	f, err := os.Open("src/config/database_config.yml")
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&Cfg.DB)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	initGoogle()
	initDB()
}
