package aruba

// ListImages lists available images from Aruba cloud provider
func (c *Client) ListImages() ([]Image, error) {
	// TODO: Implement image listing functionality
	return []Image{}, nil
}

// Image represents a cloud image
type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// TODO: Add more image fields as needed
}