package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds the entire configuration structure.
type Config struct {
	DefaultProfile string             `mapstructure:"default_profile"`
	Profiles       map[string]Profile `mapstructure:"profiles"`
}

// Profile contains the configuration for a single provider profile.
type Profile struct {
	Provider string                 `mapstructure:"provider"`
	Settings map[string]interface{} `mapstructure:"settings"`
}

// ArubaSettings defines the specific configuration for the Aruba provider.
type ArubaSettings struct {
	Datacenter   string `mapstructure:"datacenter"`
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

var vp *viper.Viper

// Load initializes and loads the configuration from file.
func Load(configPath string) (*Config, error) {
	vp = viper.New()

	if configPath != "" {
		vp.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		vp.AddConfigPath(fmt.Sprintf("%s/.seca", home))
		vp.SetConfigName("config")
		vp.SetConfigType("yaml")
	}

	// Set defaults
	vp.SetDefault("default_profile", "default")

	if err := vp.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error and use defaults or env vars
			// In a real app, you might want to prompt the user to create one.
			fmt.Println("Configuration file not found. Using defaults.")
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := vp.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetArubaSettings extracts Aruba-specific settings from a profile.
func (p *Profile) GetArubaSettings() (*ArubaSettings, error) {
	if p.Provider != "aruba" {
		return nil, fmt.Errorf("profile is not for aruba provider, but for '%s'", p.Provider)
	}

	// A bit of manual mapping to get the structure right
	settings := &ArubaSettings{}
	if dc, ok := p.Settings["datacenter"].(string); ok {
		settings.Datacenter = dc
	}
	if clientID, ok := p.Settings["client_id"].(string); ok {
		settings.ClientID = clientID
	}
	if clientSecret, ok := p.Settings["client_secret"].(string); ok {
		settings.ClientSecret = clientSecret
	}

	return settings, nil
}
