package aruba

import (
	"fmt"
	"net/http"
)

// Client represents an Aruba cloud provider client
type Client struct {
	baseURL    string
	datacenter string
	httpClient *http.Client
	username   string
	password   string
}

// ClientConfig holds the configuration for the Aruba client.
type ClientConfig struct {
	Datacenter string
	Username   string
	Password   string
}

const arubaAPIBaseURL = "https://api.arubacloud.com"

// NewClient creates a new Aruba client.
func NewClient(config ClientConfig) (*Client, error) {
	if config.Datacenter == "" {
		return nil, fmt.Errorf("Datacenter is required")
	}
	return &Client{
		baseURL:    arubaAPIBaseURL, // Use the correct, fixed base URL
		datacenter: config.Datacenter,
		httpClient: &http.Client{},
		username:   config.Username,
		password:   config.Password,
	}, nil
}