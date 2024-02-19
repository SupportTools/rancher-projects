package rancher

import (
	"fmt"
	"log"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// Project processes a project within a specified cluster, verifying the project and namespace, and optionally generating a kubeconfig.
func Project(cfg *config.Config, clusterID string) error {
	// Assuming VerifyProject has been refactored to return an error.
	if err := VerifyProject(cfg, clusterID, cfg.ProjectName); err != nil {
		log.Printf("Error verifying project '%s': %v\n", cfg.ProjectName, err)
		// Decide whether to return or continue based on your error handling policy.
		return fmt.Errorf("error verifying project '%s': %v", cfg.ProjectName, err)
	}

	if cfg.ProjectName != "" {
		// Assuming VerifyNamespace has been refactored to return an error.
		if err := VerifyNamespace(cfg, clusterID, cfg.ProjectName); err != nil {
			log.Printf("Error verifying namespace '%s': %v\n", cfg.ProjectName, err)
			// Decide whether to return or continue based on your error handling policy.
			return fmt.Errorf("error verifying namespace '%s': %v", cfg.ProjectName, err)
		}
	}

	if cfg.CreateKubeconfig {
		// Assuming GenerateKubeconfig has been refactored to return an error.
		if err := GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID); err != nil {
			log.Printf("Error generating kubeconfig for cluster '%s': %v\n", clusterID, err)
			// Decide whether to return or continue based on your error handling policy.
			return fmt.Errorf("error generating kubeconfig for cluster '%s': %v", clusterID, err)
		}
	}
	return nil
}
