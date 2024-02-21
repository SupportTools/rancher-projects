package rancher

import (
	"fmt"
	"log"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// MainProject processes a project within a specified cluster, verifying the project and namespace, and optionally generating a kubeconfig.
func MainProject(cfg *config.Config, clusterID string) error {
	fmt.Println("Verifying project...")
	if err := VerifyProject(cfg, clusterID, cfg.ProjectName); err != nil {
		log.Printf("Error verifying project '%s': %v\n", cfg.ProjectName, err)
		return fmt.Errorf("error verifying project '%s': %v", cfg.ProjectName, err)
	}

	fmt.Println("Verifying namespace...")
	if cfg.ProjectName != "" {
		if err := VerifyNamespace(cfg, clusterID, cfg.Namespace); err != nil {
			log.Printf("Error verifying namespace '%s': %v\n", cfg.Namespace, err)
			return fmt.Errorf("error verifying namespace '%s': %v", cfg.Namespace, err)
		}
	}

	fmt.Printf("Creating kubeconfig for cluster '%s'...\n", clusterID)
	if cfg.CreateKubeconfig {
		if err := GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID); err != nil {
			log.Printf("Error generating kubeconfig for cluster '%s': %v\n", clusterID, err)
			return fmt.Errorf("error generating kubeconfig for cluster '%s': %v", clusterID, err)
		}
	}
	return nil
}
