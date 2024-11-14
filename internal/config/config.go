package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/yardnsm/leeches/pkg/parcel"
)

const (
	// Hardcoded expiration time for requests - 1hr.
	MaxChargeRequestTime = 1 * time.Hour
)

type LeechesConfig struct {
	TelegramToken string `json:"telegramToken"`
	Database      string `json:"database"`

	Flavor string `json:"flavor"`

	CredentialsParcelPassword string `json:"credentialsParcelPassword"`
	CreditCardParcelPassword  string `json:"creditCardParcelPassword"`

	Webhook LeechesWebhookConfig `json:"webhook"`
}

type LeechesWebhookConfig struct {
	Port      string `json:"port"`
	PublicURL string `json:"publicUrl"`
	Cert      string `json:"certificate"`
}

type HeverCredentialsConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type HeverCreditCardConfig struct {
	Number string `json:"number"`
	Year   string `json:"year"`
	Month  string `json:"month"`
}

func LoadConfig(path string) (LeechesConfig, error) {
	var config LeechesConfig

	configJson, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("Error when opening config file: %v", err)
	}

	err = json.Unmarshal(configJson, &config)
	if err != nil {
		return config, fmt.Errorf("Error when parsing config file: %v", err)
	}

	return config, nil
}

func LoadCredentialsConfig(path string, password []byte) (HeverCredentialsConfig, error) {
	var config HeverCredentialsConfig

	credsParcel, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("Error when opening credentials parcel: %v", err)
	}

	err = parcel.Unmarshal(credsParcel, password, &config)
	if err != nil {
		return config, fmt.Errorf("Error when decrypting credentials parcel file: %v", err)
	}

	return config, nil
}

func LoadCreditCardConfig(path string, password []byte) (HeverCreditCardConfig, error) {
	var config HeverCreditCardConfig

	credsParcel, err := os.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("Error when opening credit card parcel: %v", err)
	}

	err = parcel.Unmarshal(credsParcel, password, &config)
	if err != nil {
		return config, fmt.Errorf("Error when decrypting credit card parcel file: %v", err)
	}

	return config, nil
}
