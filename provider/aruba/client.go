package aruba

import (
	"net/http"

	storageapi "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
)

// Client represents an Aruba cloud provider client
type Client struct {
	// storageClient is the SECA Storage API client for image management
	storageClient *storageapi.ClientWithResponses
	// baseURL is the base URL for the Aruba SECA API endpoint
	baseURL string
	// httpClient is the underlying HTTP client
	httpClient *http.Client
}

// ClientConfig holds configuration for the Aruba provider client
type ClientConfig struct {
	// BaseURL is the base URL for the Aruba SECA API endpoint
	BaseURL string
	// HTTPClient is the HTTP client to use (optional, defaults to http.DefaultClient)
	HTTPClient *http.Client
}

// NewClient creates a new Aruba provider client
func NewClient(config ClientConfig) (*Client, error) {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// Initialize the SECA Storage API client
	storageClient, err := storageapi.NewClientWithResponses(config.BaseURL, storageapi.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}

	return &Client{
		storageClient: storageClient,
		baseURL:       config.BaseURL,
		httpClient:    httpClient,
	}, nil
}