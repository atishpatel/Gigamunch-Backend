package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config is used to load different configurations
//between local, dev, and production environments
type Config struct {
	ClientID                   string `json:"client_id"`
	CookieName                 string `json:"cookie_name"`
	BrowserAPIKey              string `json:"browser_api_key"`
	ServerAPIKey               string `json:"server_api_key"`
	GoogleAppCredentialsPath   string `json:"google_app_credentials_path"`
	RedisSessionServerIP       string `json:"redis_session_server_ip"`
	RedisSessionServerPassword string `json:"redis_session_server_password"`
}

var (
	config *Config
)

func loadConfig() {
	filedata, err := ioutil.ReadFile("private/dev_config.json")
	if err != nil {
		log.Println("Failed to open config file in private folder.")
		log.Fatal(err)
	}
	err = json.Unmarshal(filedata, &config)
	if err != nil {
		log.Println("Failed to unmarshal config file.")
		log.Fatal(err)
	}

}
