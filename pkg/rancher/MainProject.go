package rancher

import (
	"fmt"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// MainProject processes a project within a specified cluster, verifying the project and namespace, and optionally generating a kubeconfig.
func MainProject(cfg *config.Config, clusterID string) error {
	logger.Info("Starting project verification...")

	logger.Info(fmt.Sprintf("Verifying project: %s", cfg.ProjectName))
	if err := VerifyProject(cfg, clusterID, cfg.ProjectName); err != nil {
		logger.Error(fmt.Sprintf("Error verifying project '%s': %v", cfg.ProjectName, err))
		return fmt.Errorf("error verifying project '%s': %v", cfg.ProjectName, err)
	}

	logger.Info(fmt.Sprintf("Verifying namespace: %s", cfg.Namespace))
	if cfg.ProjectName != "" {
		if err := VerifyNamespace(cfg, clusterID, cfg.Namespace); err != nil {
			logger.Error(fmt.Sprintf("Error verifying namespace '%s': %v", cfg.Namespace, err))
			return fmt.Errorf("error verifying namespace '%s': %v", cfg.Namespace, err)
		}
	}

	logger.Info(fmt.Sprintf("Creating kubeconfig for cluster '%s'...", clusterID))
	if cfg.CreateKubeconfig {
		if err := GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID); err != nil {
			logger.Error(fmt.Sprintf("Error generating kubeconfig for cluster '%s': %v", clusterID, err))
			return fmt.Errorf("error generating kubeconfig for cluster '%s': %v", clusterID, err)
		}
	}

	return nil
}
