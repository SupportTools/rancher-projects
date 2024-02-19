package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// MultiCluster processes multiple clusters based on configuration settings.
func MultiCluster(cfg *config.Config) error {
	if err := GetAllClusterIDs(cfg); err != nil {
		fmt.Printf("Failed to get all cluster IDs: %v\n", err)
		return fmt.Errorf("failed to get all cluster IDs: %v", err)
	}

	keyPairs := strings.Split(cfg.ClusterLabels, ",")

	for _, clusterID := range cfg.ClusterIDs {
		clusterName, _, err := ParseClusterID(clusterID)
		if err != nil {
			fmt.Printf("Error parsing cluster ID '%s': %v\n", clusterID, err)
			continue // Skip this cluster and move to the next.
		}

		// Assuming IsClusterActive has been refactored to return an error.
		active, err := IsClusterActive(clusterName, cfg)
		if err != nil {
			fmt.Printf("Failed to check if cluster %s is active: %v\n", clusterName, err)
			continue // Skip this cluster and move to the next.
		}

		if active {
			if cfg.ClusterType != "" {
				ClusterByType(clusterName, cfg)
			} else if cfg.ClusterLabels != "" {
				ClusterByLabels(clusterName, keyPairs, cfg)
			}
		} else {
			fmt.Printf("Skipping cluster %s because it is not active\n", clusterName)
		}

	}

	return nil
}
