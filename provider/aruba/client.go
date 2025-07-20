package aruba

import (
	"fmt"
	"net/http"
)

// Client represents an Aruba cloud provider client
type Client struct {
	// baseURL is the base URL for the Aruba Cloud API endpoint
	baseURL string
	// username for Aruba Cloud API authentication
	username string
	// password for Aruba Cloud API authentication
	password string
	// httpClient is the underlying HTTP client
	httpClient *http.Client
}

// ClientConfig holds configuration for the Aruba provider client
type ClientConfig struct {
	// BaseURL is the base URL for the Aruba Cloud API endpoint (e.g., "https://api.dc3.com/1.0")
	BaseURL string
	// Username for Aruba Cloud API authentication
	Username string
	// Password for Aruba Cloud API authentication
	Password string
	// HTTPClient is the HTTP client to use (optional, defaults to http.DefaultClient)
	HTTPClient *http.Client
}

// NewClient creates a new Aruba provider client
func NewClient(config ClientConfig) (*Client, error) {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("BaseURL is required (e.g., https://api.dc3.com/1.0)")
	}
	if config.Username == "" {
		return nil, fmt.Errorf("Username is required for Aruba Cloud API")
	}
	if config.Password == "" {
		return nil, fmt.Errorf("Password is required for Aruba Cloud API")
	}

	return &Client{
		baseURL:    config.BaseURL,
		username:   config.Username,
		password:   config.Password,
		httpClient: httpClient,
	}, nil
}