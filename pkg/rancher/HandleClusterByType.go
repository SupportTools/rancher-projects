package rancher

import (
	"fmt"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func HandleClusterByType(clusterName, clusterID string, cfg *config.Config) {
	fmt.Printf("Handling cluster %s by type...\n", clusterName)

	clusterType := GetClusterType(cfg, clusterName)
	if clusterType != "" {
		fmt.Printf("Cluster type: %s\n", clusterType)

		fmt.Println("Performing actions based on cluster type...")
	} else {
		fmt.Printf("Failed to retrieve cluster type for %s\n", clusterName)
	}
}
