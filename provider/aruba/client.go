package aruba

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// ArubaProject represents a project from the Aruba API
type ArubaProject struct {
	ID         string            `json:"id"`
	Properties ProjectProperties `json:"properties"`
}

// ProjectProperties contains the properties of a project
type ProjectProperties struct {
	Default bool `json:"default"`
}

// GetDefaultProjectID fetches all projects and returns the ID of the default one.
func (c *Client) GetDefaultProjectID(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/projects", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for projects: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Aruba API for projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Aruba API returned status %d for projects: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read projects response body: %w", err)
	}

	log.Printf("Aruba Projects Response: %s", string(body))

	var projects []ArubaProject
	if err := json.Unmarshal(body, &projects); err != nil {
		return "", fmt.Errorf("failed to decode projects response: %w", err)
	}

	for _, p := range projects {
		if p.Properties.Default {
			return p.ID, nil
		}
	}

	return "", fmt.Errorf("no default project found")
}