package config

// Package config is a singleton pattern design class that allows different components
// of the server to find the configurations in current environment
//
// This class is useful because environment variables change for different stages
// of the development process from: local, dev, prod

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/appengine"
)

const (
	kindConfig = "ServerConfig"
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

// GetConfig returns the configurations of the server
func GetConfig() *Config {
	if config == nil {
		loadConfig()
	}
	return config

}

func loadConfig() {
	if appengine.IsDevAppServer() {
		privateDirPath := os.Getenv("GIGAMUNCH_PRIVATE_DIR")
		if privateDirPath == "" {
			log.Fatal("environment variable GIGAMUNCH_PRIVATE_DIR not set")
		}
		filedata, err := ioutil.ReadFile(privateDirPath + "/dev_config.json")
		if err != nil {
			log.Println("Failed to open config file in private folder.")
			log.Fatal(err)
		}
		err = json.Unmarshal(filedata, &config)
		if err != nil {
			log.Println("Failed to unmarshal config file.")
			log.Fatal(err)
		}
	} else {
		// TODO(Atish): load from metadata on project
	}
}

func init() {
	loadConfig()
}
