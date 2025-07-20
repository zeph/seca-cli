package aruba

import (
	"context"
	"fmt"

	storageapi "github.com/eu-sovereign-cloud/go-sdk/pkg/spec/foundation.storage.v1"
)

// ListImages lists available images from Aruba cloud provider using SECA Storage API
func (c *Client) ListImages(ctx context.Context, tenant string) (*storageapi.ImageIterator, error) {
	if c.storageClient == nil {
		return nil, fmt.Errorf("storage client not initialized")
	}

	// Use the generated SECA Storage API client to list images
	resp, err := c.storageClient.ListImagesWithResponse(ctx, tenant, &storageapi.ListImagesParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode())
	}

	return resp.JSON200, nil
}

// GetImage retrieves a specific image by name from Aruba cloud provider
func (c *Client) GetImage(ctx context.Context, tenant, name string) (*storageapi.Image, error) {
	if c.storageClient == nil {
		return nil, fmt.Errorf("storage client not initialized")
	}

	resp, err := c.storageClient.GetImageWithResponse(ctx, tenant, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get image %s: %w", name, err)
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("API returned status %d for image %s", resp.StatusCode(), name)
	}

	return resp.JSON200, nil
}