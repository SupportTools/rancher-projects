package rancher

import (
	"fmt"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// ClusterByType performs actions based on the type of the specified cluster.
func ClusterByType(clusterName string, cfg *config.Config) {
	fmt.Printf("Handling cluster %s by type...\n", clusterName)

	// Assuming GetClusterType has been refactored to return an error.
	clusterType, err := GetClusterType(cfg, clusterName)
	if err != nil {
		fmt.Printf("Failed to retrieve cluster type for %s: %v\n", clusterName, err)
		return // Early return on error.
	}

	if clusterType != "" {
		fmt.Printf("Cluster type: %s\n", clusterType)
		fmt.Println("Performing actions based on cluster type...")
	} else {
		fmt.Printf("Cluster type for %s is undefined or empty\n", clusterName)
	}
}
