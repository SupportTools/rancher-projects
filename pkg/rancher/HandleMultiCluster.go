package rancher

import (
	"fmt"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func HandleMultiCluster(cfg *config.Config) {
	// Get the cluster IDs and populate the ClusterIDs field in the config
	GetAllClusterIDs(cfg)

	keyPairs := strings.Split(cfg.ClusterLabels, ",")
	keyPairCount := len(keyPairs)

	for _, clusterID := range cfg.ClusterIDs {
		clusterName, parsedClusterID := ParseClusterID(clusterID)

		if IsClusterActive(clusterName, cfg) {
			if cfg.ClusterType != "" {
				HandleClusterByType(clusterName, parsedClusterID, cfg)
			} else if cfg.ClusterLabels != "" {
				HandleClusterByLabels(clusterName, parsedClusterID, keyPairs, keyPairCount, cfg)
			}
		} else {
			fmt.Println("Skipping Cluster because it is not Active")
		}

		fmt.Println("##################################################################################")
	}
}
