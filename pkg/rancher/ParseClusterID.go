package rancher

import (
	"fmt"
	"strings"
)

// ParseClusterID attempts to parse the cluster name and ID from a given string.
// It returns an error if the input format is not as expected.
func ParseClusterID(clusterID string) (clusterName, parsedClusterID string, err error) {
	parts := strings.Split(clusterID, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid clusterID format: %s", clusterID)
	}
	clusterName, parsedClusterID = parts[0], parts[1]
	return clusterName, parsedClusterID, nil
}
