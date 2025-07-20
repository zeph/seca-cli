package aruba

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const arubaAuthURL = "https://login.aruba.it/auth/realms/cmp-new-apikey/protocol/openid-connect/token"

// Client represents an Aruba cloud provider client
type Client struct {
	baseURL     string
	datacenter  string
	httpClient  *http.Client
	accessToken string
}

// ClientConfig holds the configuration for the Aruba client.
type ClientConfig struct {
	Datacenter   string
	ClientID     string
	ClientSecret string
}

// TokenResponse models the response from the Aruba authentication server.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// NewClient creates a new Aruba client and authenticates to get a token.
func NewClient(config ClientConfig) (*Client, error) {
	if config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("client_id and client_secret are required")
	}

	client := &Client{
		baseURL:    "https://api.arubacloud.com/v2",
		datacenter: config.Datacenter,
		httpClient: &http.Client{},
	}

	// Get the access token upon creation
	accessToken, err := client.getAccessToken(config.ClientID, config.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Aruba: %w", err)
	}
	client.accessToken = accessToken

	return client, nil
}

// getAccessToken fetches an OAuth2 token from Aruba's auth endpoint.
func (c *Client) getAccessToken(clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", arubaAuthURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create token request: %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send token request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read token response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Aruba auth API returned status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	log.Println("Successfully obtained Aruba access token.")
	return tokenResponse.AccessToken, nil
}

// ProjectsResponse is the top-level response for the /projects endpoint
type ProjectsResponse struct {
	Total  int            `json:"total"`
	Values []ArubaProject `json:"values"`
}

// ArubaProject represents a project from the Aruba API
type ArubaProject struct {
	Metadata   ProjectMetadata   `json:"metadata"`
	Properties ProjectProperties `json:"properties"`
}

// ProjectMetadata contains the metadata of a project
type ProjectMetadata struct {
	ID string `json:"id"`
}

// ProjectProperties contains the properties of a project
type ProjectProperties struct {
	Default bool `json:"default"`
}

// GetDefaultAccountID fetches all accounts and returns the ID of the first one as a stand-in for projects.
func (c *Client) GetDefaultAccountID(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://api.arubacloud.com/v2/accounts", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for accounts: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Aruba API for accounts: %w", err)
		return "", fmt.Errorf("failed to call Aruba API for projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Aruba API returned status %d for accounts: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read accounts response body: %w", err)
	}

	log.Printf("Aruba Projects Response: %s", string(body))

	// For now, we'll assume the account response is a simple array of objects with an ID.
	// This is a placeholder until the real account structure is known.
	var accounts []map[string]interface{}
	if err := json.Unmarshal(body, &accounts); err != nil {
		return "", fmt.Errorf("failed to decode accounts response: %w", err)
	}

	if len(accounts) > 0 && accounts[0]["id"] != nil {
		if id, ok := accounts[0]["id"].(string); ok {
			return id, nil
		}
	}

	return "", fmt.Errorf("no accounts found or first account has no ID")
}