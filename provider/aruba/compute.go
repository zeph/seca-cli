package aruba

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	storageapi "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
)

// ArubaOSImage represents an OS image from Aruba Cloud API
type ArubaOSImage struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Architecture string `json:"architecture"`
}

// ListImages lists available OS images from Aruba Cloud API and converts them to SECA format
func (c *Client) ListImages(ctx context.Context, tenant string) (*storageapi.ImageIterator, error) {
	// Call real Aruba Cloud API to get OS images
	arubaImages, err := c.getArubaOSImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Aruba OS images: %w", err)
	}

	// Convert Aruba images to SECA format
	secaImages := make([]storageapi.Image, 0, len(arubaImages))
	for _, arubaImg := range arubaImages {
		secaImg := c.convertArubaImageToSECA(arubaImg, tenant)
		secaImages = append(secaImages, secaImg)
	}

	return &storageapi.ImageIterator{
		Items: secaImages,
	}, nil
}

// getArubaOSImages calls the real Aruba Cloud API to get OS images
func (c *Client) getArubaOSImages(ctx context.Context) ([]ArubaOSImage, error) {
	// Build the API URL for OS dictionary
	apiURL, err := url.JoinPath(c.baseURL, "dictionary", "os")
	if err != nil {
		return nil, fmt.Errorf("failed to build API URL: %w", err)
	}

	// Create HTTP request with basic auth
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")

	// Make the API call
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Aruba API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Aruba API returned status %d", resp.StatusCode)
	}

	// Parse the response
	var images []ArubaOSImage
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return images, nil
}

// convertArubaImageToSECA converts an Aruba OS image to SECA format
func (c *Client) convertArubaImageToSECA(arubaImg ArubaOSImage, tenant string) storageapi.Image {
	// Convert architecture (Aruba uses different naming)
	cpuArch := storageapi.Amd64 // Default
	if arubaImg.Architecture == "x86_64" || arubaImg.Architecture == "amd64" {
		cpuArch = storageapi.Amd64
	} else if arubaImg.Architecture == "arm64" || arubaImg.Architecture == "aarch64" {
		cpuArch = storageapi.Arm64
	}

	// Create SECA-compliant image
	imageName := fmt.Sprintf("aruba-os-%d", arubaImg.ID)
	boot := storageapi.UEFI
	initializer := storageapi.None
	state := storageapi.ResourceStateActive

	// Create labels and annotations
	labels := map[string]string{
		"provider":     "aruba",
		"os":          "linux", // Assuming Linux for now
		"aruba-id":    fmt.Sprintf("%d", arubaImg.ID),
	}
	annotations := map[string]string{
		"name":        arubaImg.Name,
		"description": arubaImg.Description,
		"source":      "aruba-cloud",
	}

	return storageapi.Image{
		Labels:      &labels,
		Annotations: &annotations,
		Metadata: &storageapi.RegionalResourceMetadata{
			Name:       imageName,
			ApiVersion: "storage.seca.cloud/v1",
			Kind:       storageapi.RegionalResourceMetadataKindImage,
			Tenant:     tenant,
		},
		Spec: storageapi.ImageSpec{
			BlockStorageRef: fmt.Sprintf("aruba-images/%d", arubaImg.ID),
			CpuArchitecture: cpuArch,
			Boot:            &boot,
			Initializer:     &initializer,
		},
		Status: &storageapi.ImageStatus{
			State: &state,
		},
	}
}

// GetImage retrieves a specific image by name from Aruba cloud provider
func (c *Client) GetImage(ctx context.Context, tenant, name string) (*storageapi.Image, error) {
	// Get all images and find the one with matching name
	images, err := c.ListImages(ctx, tenant)
	if err != nil {
		return nil, err
	}

	for _, img := range images.Items {
		if img.Metadata.Name == name {
			return &img, nil
		}
	}

	return nil, fmt.Errorf("image %s not found", name)
}