package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// MultiCluster processes multiple clusters based on configuration settings.
func MultiCluster(cfg *config.Config) error {
	logger.Info("Fetching all cluster IDs...")

	if err := GetAllClusterIDs(cfg); err != nil {
		logger.Error(fmt.Sprintf("Failed to get all cluster IDs: %v", err))
		return fmt.Errorf("failed to get all cluster IDs: %v", err)
	}

	keyPairs := strings.Split(cfg.ClusterLabels, ",")
	logger.Debug(fmt.Sprintf("Parsed cluster labels: %v", keyPairs))

	for _, clusterID := range cfg.ClusterIDs {
		clusterName, _, err := ParseClusterID(clusterID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error parsing cluster ID '%s': %v", clusterID, err))
			continue // Skip this cluster and move to the next.
		}

		logger.Info(fmt.Sprintf("Checking if cluster %s is active...", clusterName))
		active, err := IsClusterActive(clusterName, cfg)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to check if cluster %s is active: %v", clusterName, err))
			continue // Skip this cluster and move to the next.
		}

		if active {
			logger.Info(fmt.Sprintf("Processing active cluster: %s", clusterName))
			if cfg.ClusterType != "" {
				logger.Info(fmt.Sprintf("Processing cluster %s by type...", clusterName))
				ClusterByType(clusterName, cfg)
			} else if cfg.ClusterLabels != "" {
				logger.Info(fmt.Sprintf("Processing cluster %s by labels...", clusterName))
				ClusterByLabels(clusterName, keyPairs, cfg)
			}
		} else {
			logger.Warn(fmt.Sprintf("Skipping cluster %s because it is not active", clusterName))
		}
	}
	return nil
}
