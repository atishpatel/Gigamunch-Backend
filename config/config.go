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

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
)

// Config is the configuration loaded from datastore
type Config struct {
	JWTSecret                 string   `json:"jwt_secret" datastore:",noindex"`
	ClientID                  string   `json:"client_id" datastore:",noindex"`
	ServerKey                 string   `json:"server_key" datastore:",noindex"`
	BTEnvironment             string   `json:"bt_environment" datastore:",noindex"`
	BTMerchantID              string   `json:"bt_merchant_id" datastore:",noindex"`
	BTPublicKey               string   `json:"bt_public_key" datastore:",noindex"`
	BTPrivateKey              string   `json:"bt_private_key" datastore:",noindex"`
	TwilioAccountSID          string   `json:"twilio_account_sid" datastore:",noindex"`
	TwilioKeySID              string   `json:"twilio_key_sid" datastore:",noindex"`
	TwilioAuthToken           string   `json:"twilio_auth_token" datastore:",noindex"`
	TwilioKeyAuthToken        string   `json:"twilio_key_auth_token" datastore:",noindex"`
	TwilioIPMessagingSID      string   `json:"twilio_ip_messaging_sid" datastore:",noindex"`
	TwilioPushCredentialSID   string   `json:"push_credential_sid" datastore:",noindex"`
	TwilioAdminServiceSID     string   `json:"twilio_admin_service_sid" datastore:",noindex"`
	TwilioDeliveryServiceSID  string   `json:"twilio_delivery_service_sid" datastore:",noindex"`
	TwilioBagServiceSID       string   `json:"twilio_bag_service_sid" datastore:",noindex"`
	PhoneNumbers              []string `json:"phone_numbers" datastore:",noindex"`
	BucketName                string   `json:"bucket_name" datastore:",noindex"`
	ProjectID                 string   `json:"project_id" datastore:",noindex"`
	CompanyCardCardholderName string   `json:"company_card_cardholder_name" datastore:",noindex"`
	CompanyCardNumber         string   `json:"company_card_number" datastore:",noindex"`
	CompanyCardExpirationDate string   `json:"company_card_expiration_date" datastore:",noindex"`
	CompanyCardCVV            string   `json:"company_card_cvv" datastore:",noindex"`
	SendGridKey               string   `json:"send_grid_key" datastore:",noindex"`
	DripAPIKey                string   `json:"drip_api_key"`
	DripAccountID             string   `json:"drip_account_id"`
}

// BTEnvironment is the environment type for braintree
type BTEnvironment string

const (
	// BTSandbox is the Braintree sandbox env
	BTSandbox = "sandbox"
	// BTProduction is the Braintree production env
	BTProduction   = "production"
	privateDirPath = "../private"
)

// BTConfig has all the config needed for Braintree
type BTConfig struct {
	BTEnvironment string     `json:"bt_environment"`
	BTMerchantID  string     `json:"bt_merchant_id"`
	BTPublicKey   string     `json:"bt_public_key"`
	BTPrivateKey  string     `json:"bt_private_key"`
	CompanyCard   CreditCard `json:"company_card"`
}

// GitkitConfig is used to load different configurations
// between local, dev, and production environments
type GitkitConfig struct {
	JWTSecret                string `json:"jwt_secret"`
	ClientID                 string `json:"client_id"`
	BrowserAPIKey            string `json:"browser_api_key"`
	ServerAPIKey             string `json:"server_api_key"`
	GoogleAppCredentialsPath string `json:"google_app_credentials_path"`
}

// TwilioConfig is used to load twilio configurations
type TwilioConfig struct {
	AccountSID         string   `json:"account_sid"`
	KeySID             string   `json:"key_sid"`
	AuthToken          string   `json:"auth_token"`
	KeyAuthToken       string   `json:"key_auth_token"`
	IPMessagingSID     string   `json:"ip_messaging_sid"`
	PushCredentialSID  string   `json:"push_credential_sid"`
	PhoneNumbers       []string `json:"phone_numbers"`
	AdminServiceSID    string   `json:"admin_service_sid"`
	DeliveryServiceSID string   `json:"delivery_service_sid"`
	BagServiceSID      string   `json:"bag_service_sid"`
}

// CreditCard holds data related to a credit card.
type CreditCard struct {
	CardholderName string `json:"cardholder_name"`
	Number         string `json:"number"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}

type MailConfig struct {
	SendGridKey   string `json:"send_grid_key"`
	DripAPIKey    string `json:"drip_api_key"`
	DripAccountID string `json:"drip_account_id"`
}

var (
	gitkitConfig *GitkitConfig
	config       *Config
)

func GetConfig(ctx context.Context) Config {
	var config Config
	if appengine.IsDevAppServer() {
		filedata := readFile("config.json")
		err := json.Unmarshal(filedata, &config)
		if err != nil {
			log.Println("Failed to unmarshal config file.")
			log.Fatal(err)
		}
	} else {
		key := datastore.NewKey(ctx, "Config", "", 100, nil)
		err := datastore.Get(ctx, key, &config)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				config.PhoneNumbers = []string{"14243484448"}
				_, _ = datastore.Put(ctx, key, config)
			} else {
				log.Fatalf("Error getting Config from datastore: %+v", err)
			}
		}
	}
	return config
}

func GetMailConfig(ctx context.Context) MailConfig {
	var mailConfig MailConfig
	if appengine.IsDevAppServer() {
		filedata := readFile("mail_config.json")
		err := json.Unmarshal(filedata, &mailConfig)
		if err != nil {
			log.Println("Failed to unmarshal mail_config file.")
			log.Fatal(err)
		}
	} else {
		config := getDatastoreConfig(ctx)
		mailConfig.SendGridKey = config.SendGridKey
		mailConfig.DripAPIKey = config.DripAPIKey
		mailConfig.DripAccountID = config.DripAccountID
	}
	return mailConfig
}

// GetTwilioConfig returns the twilio configs
func GetTwilioConfig(ctx context.Context) TwilioConfig {
	var twilioConfig TwilioConfig
	if appengine.IsDevAppServer() {
		filedata := readFile("twilio_config.json")
		err := json.Unmarshal(filedata, &twilioConfig)
		if err != nil {
			log.Println("Failed to unmarshal twilio_config file.")
			log.Fatal(err)
		}
	} else {
		getDatastoreConfig(ctx)
		twilioConfig.AccountSID = config.TwilioAccountSID
		twilioConfig.AuthToken = config.TwilioAuthToken
		twilioConfig.KeySID = config.TwilioKeySID
		twilioConfig.KeyAuthToken = config.TwilioKeyAuthToken
		twilioConfig.IPMessagingSID = config.TwilioIPMessagingSID
		twilioConfig.PushCredentialSID = config.TwilioPushCredentialSID
		twilioConfig.PhoneNumbers = config.PhoneNumbers
		twilioConfig.AdminServiceSID = config.TwilioAdminServiceSID
		twilioConfig.DeliveryServiceSID = config.TwilioDeliveryServiceSID
		twilioConfig.BagServiceSID = config.TwilioBagServiceSID
	}
	return twilioConfig
}

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
		btConfig.CompanyCard.CardholderName = config.CompanyCardCardholderName
		btConfig.CompanyCard.Number = config.CompanyCardNumber
		btConfig.CompanyCard.CVV = config.CompanyCardCVV
		btConfig.CompanyCard.ExpirationDate = config.CompanyCardExpirationDate
	}
	return btConfig
}

// GetBucketName returns the bucket for uploading images
func GetBucketName(ctx context.Context) string {
	if appengine.IsDevAppServer() {
		return "gigamunch-dev-images"
	}
	getDatastoreConfig(ctx)
	return config.BucketName
}

// GetProjectID returns the project id
func GetProjectID(ctx context.Context) string {
	if appengine.IsDevAppServer() {
		return "gigamunch-omninexus-dev"
	}
	getDatastoreConfig(ctx)
	return config.ProjectID
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
	if gitkitConfig == nil || gitkitConfig.JWTSecret == "" {
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
			log.Println("Failed to unmarshal gitkit config file from.")
			log.Fatal(err)
		}
	} else {
		getDatastoreConfig(ctx)
		gitkitConfig = &GitkitConfig{
			JWTSecret: config.JWTSecret,
			ClientID:  config.ClientID,
		}
	}
}

func getDatastoreConfig(ctx context.Context) *Config {
	if config == nil || config.JWTSecret == "" {
		configTmp := new(Config)
		key := datastore.NewKey(ctx, "Config", "", 100, nil)
		err := datastore.Get(ctx, key, configTmp)
		if err != nil {
			if err == datastore.ErrNoSuchEntity {
				configTmp.PhoneNumbers = []string{"14243484448"}
				_, _ = datastore.Put(ctx, key, configTmp)
			} else {
				log.Fatalf("Error getting Config from datastore: %+v", err)
			}
		}
		config = configTmp
	}
	return config
}

func readFile(fileName string) []byte {
	filedata, err := ioutil.ReadFile(privateDirPath + "/" + fileName)
	if err != nil {
		log.Printf("Failed to open %s file in private folder(path: '%s').", fileName, privateDirPath)
		log.Fatal(err)
	}
	return filedata
}
