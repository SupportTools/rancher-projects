package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// ClusterByLabels checks if a cluster matches specified label criteria and performs actions based on that.
func ClusterByLabels(clusterName string, keyPairs []string, cfg *config.Config) {
	fmt.Printf("Handling cluster %s by labels...\n", clusterName)

	// Assuming FilterByClusterLabels has been refactored to return a bool and an error.
	matches, err := FilterByClusterLabels(clusterName, keyPairs, cfg)
	if err != nil {
		fmt.Printf("Error filtering cluster %s by labels: %v\n", clusterName, err)
		return // Early return on error to avoid proceeding with potentially invalid state.
	}

	if matches {
		fmt.Printf("Cluster matches label criteria: %s\n", strings.Join(keyPairs, ", "))
		fmt.Println("Performing actions based on cluster labels...")
	} else {
		fmt.Printf("Cluster %s does not match the specified labels\n", clusterName)
	}
}
