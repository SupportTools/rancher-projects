package rancher

import "strings"

func ParseClusterID(clusterID string) (string, string) {
	clusterName := strings.Split(clusterID, ":")[0]
	parsedClusterID := strings.Split(clusterID, ":")[1]
	return clusterName, parsedClusterID
}
