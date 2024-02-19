package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func HandleClusterByLabels(clusterName, clusterID string, keyPairs []string, keyPairCount int, cfg *config.Config) {
	fmt.Printf("Handling cluster %s by labels...\n", clusterName)

	if FilterByClusterLabels(clusterName, keyPairs, cfg) {
		fmt.Printf("Cluster matches label criteria: %s\n", strings.Join(keyPairs, ", "))
		fmt.Println("Performing actions based on cluster labels...")
	} else {
		fmt.Printf("Cluster %s does not match the specified labels\n", clusterName)
	}
}
