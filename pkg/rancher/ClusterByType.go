package rancher

import (
	"fmt"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// ClusterByType performs actions based on the type of the specified cluster.
func ClusterByType(clusterName string, cfg *config.Config) {
	logger.Info(fmt.Sprintf("Handling cluster %s by type...", clusterName))

	// Assuming GetClusterType has been refactored to return an error.
	clusterType, err := GetClusterType(cfg, clusterName)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to retrieve cluster type for %s: %v", clusterName, err))
		return // Early return on error.
	}

	if clusterType != "" {
		logger.Info(fmt.Sprintf("Cluster type: %s", clusterType))
		logger.Info("Performing actions based on cluster type...")
	} else {
		logger.Warn(fmt.Sprintf("Cluster type for %s is undefined or empty", clusterName))
	}
}
