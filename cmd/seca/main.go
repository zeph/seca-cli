package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eu-sovereign-cloud/go-sdk/internal/config"
	"github.com/eu-sovereign-cloud/go-sdk/provider/aruba"
)

func main() {
	var (
		profileName = flag.String("profile", "", "Configuration profile to use")
		command     = flag.String("cmd", "list-images", "Command to execute: list-images, get-image")
		image       = flag.String("image", "", "Image name (required for get-image command)")
		tenant      = flag.String("tenant", "default", "Tenant ID for SECA compliance")
	)
	flag.Parse()

	// Load configuration
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Determine which profile to use
	profileKey := cfg.DefaultProfile
	if *profileName != "" {
		profileKey = *profileName
	}

	profile, ok := cfg.Profiles[profileKey]
	if !ok {
		log.Fatalf("Profile '%s' not found in configuration", profileKey)
	}

	// Get Aruba-specific settings from the profile
	arubaSettings, err := profile.GetArubaSettings()
	if err != nil {
		log.Fatalf("Failed to get aruba settings from profile: %v", err)
	}

	// Initialize Aruba provider
	clientConfig := aruba.ClientConfig{
		Datacenter:   arubaSettings.Datacenter,
		ClientID:     arubaSettings.ClientID,
		ClientSecret: arubaSettings.ClientSecret,
	}

	client, err := aruba.NewClient(clientConfig)
	if err != nil {
		log.Fatalf("Failed to create Aruba client: %v", err)
	}

	ctx := context.Background()

	switch *command {
	case "list-images":
		err := listImages(ctx, client, *tenant)
		if err != nil {
			log.Fatalf("Failed to list images: %v", err)
		}

	case "get-image":
		if *image == "" {
			fmt.Fprintf(os.Stderr, "Error: -image is required for get-image command\n")
			flag.Usage()
			os.Exit(1)
		}
		err := getImage(ctx, client, *tenant, *image)
		if err != nil {
			log.Fatalf("Failed to get image: %v", err)
		}

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %s\n", *command)
		flag.Usage()
		os.Exit(1)
	}
}

func listImages(ctx context.Context, client *aruba.Client, tenant string) error {
	fmt.Printf("Listing images for tenant: %s\n", tenant)
	
	images, err := client.ListImages(ctx, tenant)
	if err != nil {
		return err
	}

	if images == nil || len(images.Items) == 0 {
		fmt.Println("No images found")
		return nil
	}

	fmt.Printf("Found %d images:\n\n", len(images.Items))
	
	for i, img := range images.Items {
		fmt.Printf("Image %d:\n", i+1)
		fmt.Printf("  Name: %s\n", img.Metadata.Name)
		
		if img.Labels != nil {
			fmt.Printf("  Labels:\n")
			for k, v := range *img.Labels {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
		
		if img.Annotations != nil {
			fmt.Printf("  Annotations:\n")
			for k, v := range *img.Annotations {
				fmt.Printf("    %s: %s\n", k, v)
			}
		}
		
		fmt.Printf("  Spec:\n")
		fmt.Printf("    CPU Architecture: %s\n", img.Spec.CpuArchitecture)
		if img.Spec.Boot != nil {
			fmt.Printf("    Boot Type: %s\n", *img.Spec.Boot)
		}
		if img.Spec.Initializer != nil {
			fmt.Printf("    Initializer: %s\n", *img.Spec.Initializer)
		}
		
		if img.Status != nil && img.Status.State != nil {
			fmt.Printf("  Status: %s\n", *img.Status.State)
		}
		
		fmt.Println()
	}

	return nil
}

func getImage(ctx context.Context, client *aruba.Client, tenant, imageName string) error {
	fmt.Printf("Getting image '%s' for tenant: %s\n", imageName, tenant)
	
	img, err := client.GetImage(ctx, tenant, imageName)
	if err != nil {
		return err
	}

	fmt.Printf("Image Details:\n")
	fmt.Printf("  Name: %s\n", img.Metadata.Name)
	
	if img.Labels != nil {
		fmt.Printf("  Labels:\n")
		for k, v := range *img.Labels {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
	
	if img.Annotations != nil {
		fmt.Printf("  Annotations:\n")
		for k, v := range *img.Annotations {
			fmt.Printf("    %s: %s\n", k, v)
		}
	}
	
	fmt.Printf("  Spec:\n")
	fmt.Printf("    CPU Architecture: %s\n", img.Spec.CpuArchitecture)
	if img.Spec.Boot != nil {
		fmt.Printf("    Boot Type: %s\n", *img.Spec.Boot)
	}
	if img.Spec.Initializer != nil {
		fmt.Printf("    Initializer: %s\n", *img.Spec.Initializer)
	}
	fmt.Printf("    Block Storage Ref: %v\n", img.Spec.BlockStorageRef)
	
	if img.Status != nil {
		fmt.Printf("  Status:\n")
		if img.Status.State != nil {
			fmt.Printf("    State: %s\n", *img.Status.State)
		}
		if img.Status.Conditions != nil {
			fmt.Printf("    Conditions:\n")
			for _, cond := range img.Status.Conditions {
				fmt.Printf("      - Type: %s, State: %s\n", 
					getStringValue(cond.Type), cond.State)
				if cond.Reason != nil {
					fmt.Printf("        Reason: %s\n", *cond.Reason)
				}
				if cond.Message != nil {
					fmt.Printf("        Message: %s\n", *cond.Message)
				}
			}
		}
	}

	return nil
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
