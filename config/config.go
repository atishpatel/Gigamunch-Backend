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

	"golang.org/x/net/context"

	"github.com/atishpatel/Gigamunch-Backend/utils"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Config is the configuration loaded from datastore
type Config struct {
	JWTSecret     string
	ClientID      string
	ServerKey     string
	BTEnvironment string
	BTMerchantID  string
	BTPublicKey   string
	BTPrivateKey  string
}

// BTEnvironment is the environment type for braintree
type BTEnvironment string

const (
	// BTSandbox is the Braintree sandbox env
	BTSandbox = "sandbox"
	// BTProduction is the Braintree production env
	BTProduction = "production"
)

// BTConfig has all the config needed for Braintree
type BTConfig struct {
	BTEnvironment string `json:"bt_environment"`
	BTMerchantID  string `json:"bt_merchant_id"`
	BTPublicKey   string `json:"bt_public_key"`
	BTPrivateKey  string `json:"bt_private_key"`
}

// GitkitConfig is used to load different configurations
//between local, dev, and production environments
type GitkitConfig struct {
	JWTSecret                string `json:"jwt_secret"`
	ClientID                 string `json:"client_id"`
	BrowserAPIKey            string `json:"browser_api_key"`
	ServerAPIKey             string `json:"server_api_key"`
	GoogleAppCredentialsPath string `json:"google_app_credentials_path"`
}

var (
	gitkitConfig   *GitkitConfig
	config         *Config
	privateDirPath string
)

// GetBTConfig returns the Braintree config
func GetBTConfig(ctx context.Context) BTConfig {
	var btConfig BTConfig
	if appengine.IsDevAppServer() {
		filedata := readFile("bt_config.json")
		err := json.Unmarshal(filedata, &btConfig)
		if err != nil {
			log.Println("Failed to unmarshal bt_config file.")
			log.Fatal(err)
		}
	} else {
		getDatastoreConfig(ctx)
		btConfig.BTEnvironment = config.BTEnvironment
		btConfig.BTMerchantID = config.BTMerchantID
		btConfig.BTPublicKey = config.BTPublicKey
		btConfig.BTPrivateKey = config.BTPrivateKey
	}
	return btConfig
}

// GetServerKey returns the server key
func GetServerKey(ctx context.Context) string {
	if gitkitConfig == nil {
		loadGitkitConfig(ctx)
	}
	if appengine.IsDevAppServer() {
		return gitkitConfig.ServerAPIKey
	}
	return config.ServerKey
}

// GetGitkitConfig returns the configurations for gitkit on the server
func GetGitkitConfig(ctx context.Context) *GitkitConfig {
	if gitkitConfig == nil {
		loadGitkitConfig(ctx)
	}
	return gitkitConfig
}

func loadGitkitConfig(ctx context.Context) {
	var err error
	if appengine.IsDevAppServer() {
		filedata := readFile("gitkit_config.json")
		err = json.Unmarshal(filedata, &gitkitConfig)
		if err != nil {
			log.Println("Failed to unmarshal gitkit config file.")
			log.Fatal(err)
		}
	} else {
		getDatastoreConfig(ctx)
		gitkitConfig = &GitkitConfig{
			JWTSecret: config.JWTSecret,
			ClientID:  config.ClientID,
		}
	}

	gitkitConfig.JWTSecret, err = utils.Decrypt("KTd6M18avNkASNK149TDhyl3m45Mxqw2", gitkitConfig.JWTSecret)
	if err != nil {
		log.Fatalf("Error decoding jwt secret: %+v", err)
	}
}

func getDatastoreConfig(ctx context.Context) {
	if config == nil {
		config = new(Config)
		key := datastore.NewKey(ctx, "Config", "", 100, nil)
		err := datastore.Get(ctx, key, config)
		if err != nil {
			utils.Errorf(ctx, "getDatastoreConfig error: %+v", err)
			log.Fatalf("Error getting Config from datastore: %+v", err)
		}
	}
}

func readFile(fileName string) []byte {
	filedata, err := ioutil.ReadFile(privateDirPath + "/" + fileName)
	if err != nil {
		log.Printf("Failed to open %s file in private folder.", fileName)
		log.Fatal(err)
	}
	return filedata
}

func init() {
	if appengine.IsDevAppServer() {
		privateDirPath = os.Getenv("GIGAMUNCH_PRIVATE_DIR")
		if privateDirPath == "" {
			log.Fatal("environment variable GIGAMUNCH_PRIVATE_DIR not set")
		}
	}
}
