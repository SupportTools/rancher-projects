package rancher

import (
	"fmt"
	"log"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// SingleCluster processes a single cluster by verifying it, handling projects within it, and optionally generating a kubeconfig.
func SingleCluster(cfg *config.Config) error {

	fmt.Println("Processing single cluster...")

	fmt.Println("Verifying cluster...")
	if err := VerifyCluster(cfg); err != nil {
		log.Printf("Error verifying cluster: %v\n", err)
		return fmt.Errorf("error verifying cluster: %v", err)
	}

	fmt.Println("Getting cluster ID...")
	clusterID, err := GetClusterID(cfg)
	if err != nil {
		log.Printf("Error getting cluster ID: %v\n", err)
		return fmt.Errorf("error getting cluster ID: %v", err)
	}

	if cfg.ProjectName != "" {
		if err := MainProject(cfg, clusterID); err != nil {
			log.Printf("Error handling project '%s': %v\n", cfg.ProjectName, err)
			return fmt.Errorf("error handling project '%s': %v", cfg.ProjectName, err)
		}
	}

	return nil
}
