package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// ClusterByLabels checks if a cluster matches specified label criteria and performs actions based on that.
func ClusterByLabels(clusterName string, keyPairs []string, cfg *config.Config) {
	logger.Info(fmt.Sprintf("Handling cluster %s by labels...", clusterName))
	logger.Debug(fmt.Sprintf("Label criteria: %s", strings.Join(keyPairs, ", ")))

	// Assuming FilterByClusterLabels has been refactored to return a bool and an error.
	matches, err := FilterByClusterLabels(clusterName, keyPairs, cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Error filtering cluster %s by labels: %v", clusterName, err))
		return // Early return on error to avoid proceeding with potentially invalid state.
	}

	if matches {
		logger.Info(fmt.Sprintf("Cluster %s matches label criteria: %s", clusterName, strings.Join(keyPairs, ", ")))
		logger.Info("Performing actions based on cluster labels...")
	} else {
		logger.Info(fmt.Sprintf("Cluster %s does not match the specified labels", clusterName))
	}
}
