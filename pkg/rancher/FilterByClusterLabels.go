package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// FilterByClusterLabels checks if any of the clusters match the specified labels.
// It returns true if a match is found, false otherwise, along with an error in case of failure.
func FilterByClusterLabels(clusterName string, keyPairs []string, cfg *config.Config) (bool, error) {
	fmt.Println("Filtering by cluster label")
	fmt.Println("Cluster name:", clusterName)
	fmt.Println("Labels:", keyPairs)

	// Prepare the request to the Rancher API.
	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	req, err := http.NewRequest("GET", url, http.NoBody) // Updated to use http.NoBody
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request.
	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to filter clusters by label: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("failed to filter clusters by label. Status code: %d", resp.StatusCode)
	}

	// Read and unmarshal the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}
	var clusters []map[string]interface{}
	if err := json.Unmarshal(body, &clusters); err != nil {
		return false, fmt.Errorf("failed to parse clusters by label: %w", err)
	}

	// Iterate over the clusters and their labels to find a match.
	for _, cluster := range clusters {
		labels, ok := cluster["labels"].(map[string]interface{})
		if !ok {
			continue // Skip if labels are not in the expected format.
		}
		for _, keyPair := range keyPairs {
			parts := strings.Split(keyPair, "=")
			if len(parts) != 2 {
				return false, fmt.Errorf("invalid key-value pair: %s", keyPair)
			}

			key, value := parts[0], parts[1]
			if labelValue, exists := labels[key]; exists && labelValue == value {
				return true, nil
			}
		}
	}

	return false, nil
}
