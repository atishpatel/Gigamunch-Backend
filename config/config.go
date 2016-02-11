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

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
)

// GitkitConfig is used to load different configurations
//between local, dev, and production environments
type GitkitConfig struct {
	JWTSecret                string `json:"jwt_secret"`
	ClientID                 string `json:"client_id"`
	BrowserAPIKey            string `json:"browser_api_key"`
	ServerAPIKey             string `json:"server_api_key"`
	GoogleAppCredentialsPath string `json:"google_app_credentials_path"`
}

// SessionConfig has the configurations for storing user sessions
type SessionConfig struct {
	RedisSessionServerIP       string `json:"redis_session_server_ip"`
	RedisSessionServerPassword string `json:"redis_session_server_password"`
}

// LiveMealConfig has the configurations for querying LiveMeals
type LiveMealConfig struct {
	MySQLServerIP     string `json:"mysql_server_ip"`
	MySQLUser         string `json:"mysql_user"`
	MySQLUserPassword string `json:"mysql_user_password"`
}

var (
	gitkitConfig   *GitkitConfig
	sessionConfig  *SessionConfig
	liveMealConfig *LiveMealConfig
	privateDirPath string
)

// GetSessionConfig returns the configurations for sessions on the server
func GetSessionConfig() *SessionConfig {
	if sessionConfig == nil {
		loadSessionConfig()
	}
	return sessionConfig
}

func loadSessionConfig() {
	if appengine.IsDevAppServer() {
		filedata, err := ioutil.ReadFile(privateDirPath + "/session_config.json")
		if err != nil {
			log.Println("Failed to open config file in private folder.")
			log.Fatal(err)
		}
		err = json.Unmarshal(filedata, &sessionConfig)
		if err != nil {
			log.Println("Failed to unmarshal session config file.")
			log.Fatal(err)
		}
	} else {
		// TODO(Atish): load from metadata on project
	}
}

// GetGitkitConfig returns the configurations for gitkit on the server
func GetGitkitConfig() *GitkitConfig {
	if gitkitConfig == nil {
		loadGitkitConfig()
	}
	return gitkitConfig
}

func loadGitkitConfig() {
	if appengine.IsDevAppServer() {
		filedata, err := ioutil.ReadFile(privateDirPath + "/gitkit_config.json")
		if err != nil {
			log.Println("Failed to open config file in private folder.")
			log.Fatal(err)
		}
		err = json.Unmarshal(filedata, &gitkitConfig)
		if err != nil {
			log.Println("Failed to unmarshal gitkit config file.")
			log.Fatal(err)
		}
	} else {
		// TODO(Atish): load from metadata on project
	}
	var err error
	gitkitConfig.JWTSecret, err = utils.Decrypt("KTd6M18avNkASNK149TDhyl3m45Mxqw2", gitkitConfig.JWTSecret)
	if err != nil {
		log.Fatalf("Error decoding jwt secret: %+v", err)
	}
}

// GetLiveMealConfig returns the configurations for live meals on the server
func GetLiveMealConfig() *LiveMealConfig {
	if liveMealConfig == nil {
		loadLiveMealConfig()
	}
	return liveMealConfig
}

func loadLiveMealConfig() {
	if appengine.IsDevAppServer() {
		filedata, err := ioutil.ReadFile(privateDirPath + "/live_meal_config.json")
		if err != nil {
			log.Println("Failed to open 'live_meal_config.json' config file in private folder.")
			log.Fatal(err)
		}
		err = json.Unmarshal(filedata, &liveMealConfig)
		if err != nil {
			log.Println("Failed to unmarshal live meal config file.")
			log.Fatal(err)
		}
	} else {
		// TODO(Atish): load from metadata on project
	}
}

func init() {
	if appengine.IsDevAppServer() {
		privateDirPath = os.Getenv("GIGAMUNCH_PRIVATE_DIR")
		if privateDirPath == "" {
			log.Fatal("environment variable GIGAMUNCH_PRIVATE_DIR not set")
		}
	}
	loadGitkitConfig()
	loadSessionConfig()
	loadLiveMealConfig()
}
