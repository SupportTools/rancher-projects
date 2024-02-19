package rancher

import (
	"fmt"
	"log"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// SingleCluster processes a single cluster by verifying it, handling projects within it, and optionally generating a kubeconfig.
func SingleCluster(cfg *config.Config) error {
	// VerifyCluster now returns an error which should be handled.
	if err := VerifyCluster(cfg); err != nil {
		log.Printf("Error verifying cluster: %v\n", err)
		return fmt.Errorf("error verifying cluster: %v", err)
	}

	// GetClusterID now returns an ID and an error which should be handled.
	clusterID, err := GetClusterID(cfg)
	if err != nil {
		log.Printf("Error getting cluster ID: %v\n", err)
		return fmt.Errorf("error getting cluster ID: %v", err)
	}

	if cfg.ProjectName != "" {
		if err := Project(cfg, clusterID); err != nil {
			log.Printf("Error handling project '%s': %v\n", cfg.ProjectName, err)
			return fmt.Errorf("error handling project '%s': %v", cfg.ProjectName, err)
		}
	}

	if cfg.CreateKubeconfig {
		if err := GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID); err != nil {
			log.Printf("Error generating kubeconfig for cluster '%s': %v\n", clusterID, err)
			return fmt.Errorf("error generating kubeconfig for cluster '%s': %v", clusterID, err)
		}
	}
	return nil
}
