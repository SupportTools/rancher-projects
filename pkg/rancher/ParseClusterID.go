package rancher

import (
	"fmt"
	"strings"
)

// ParseClusterID attempts to parse the cluster name and ID from a given string.
// It returns an error if the input format is not as expected.
func ParseClusterID(clusterID string) (clusterName, parsedClusterID string, err error) {
	logger.Debug(fmt.Sprintf("Parsing cluster ID: %s", clusterID))

	parts := strings.Split(clusterID, ":")
	if len(parts) != 2 {
		logger.Error(fmt.Sprintf("Invalid clusterID format: %s", clusterID))
		return "", "", fmt.Errorf("invalid clusterID format: %s", clusterID)
	}

	clusterName, parsedClusterID = parts[0], parts[1]
	logger.Debug(fmt.Sprintf("Parsed cluster name: %s, cluster ID: %s", clusterName, parsedClusterID))

	return clusterName, parsedClusterID, nil
}
