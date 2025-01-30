package rancher

import (
	"fmt"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// SingleCluster processes a single cluster by verifying it, handling projects within it, and optionally generating a kubeconfig.
func SingleCluster(cfg *config.Config) error {
	logger.Info("Processing single cluster...")

	logger.Info("Verifying cluster...")
	if err := VerifyCluster(cfg); err != nil {
		logger.Error(fmt.Sprintf("Error verifying cluster: %v", err))
		return fmt.Errorf("error verifying cluster: %v", err)
	}

	logger.Info("Retrieving cluster ID...")
	clusterID, err := GetClusterID(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting cluster ID: %v", err))
		return fmt.Errorf("error getting cluster ID: %v", err)
	}

	if cfg.ProjectName != "" {
		logger.Info(fmt.Sprintf("Processing project '%s' in cluster '%s'...", cfg.ProjectName, clusterID))
		if err := MainProject(cfg, clusterID); err != nil {
			logger.Error(fmt.Sprintf("Error handling project '%s': %v", cfg.ProjectName, err))
			return fmt.Errorf("error handling project '%s': %v", cfg.ProjectName, err)
		}
	}

	logger.Info("Single cluster processing completed successfully.")
	return nil
}
