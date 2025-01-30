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
	logger.Info("Filtering by cluster label")
	logger.Debug(fmt.Sprintf("Cluster name: %s", clusterName))
	logger.Debug(fmt.Sprintf("Labels to match: %s", strings.Join(keyPairs, ", ")))

	// Prepare the request to the Rancher API.
	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create request: %v", err))
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute the HTTP request.
	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info("Sending GET request to fetch cluster labels...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to filter clusters by label: %v", err))
		return false, fmt.Errorf("failed to filter clusters by label: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		return false, fmt.Errorf("failed to filter clusters by label. Status code: %d", resp.StatusCode)
	}

	// Read and unmarshal the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return false, fmt.Errorf("failed to read response body: %w", err)
	}
	logger.Debug(fmt.Sprintf("Received response body: %s", string(body)))

	var clusters []map[string]interface{}
	if err := json.Unmarshal(body, &clusters); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse clusters by label: %v", err))
		return false, fmt.Errorf("failed to parse clusters by label: %w", err)
	}

	// Iterate over the clusters and their labels to find a match.
	for _, cluster := range clusters {
		labels, ok := cluster["labels"].(map[string]interface{})
		if !ok {
			logger.Warn(fmt.Sprintf("Skipping cluster %v due to missing or invalid labels format", cluster))
			continue // Skip if labels are not in the expected format.
		}
		for _, keyPair := range keyPairs {
			parts := strings.Split(keyPair, "=")
			if len(parts) != 2 {
				logger.Error(fmt.Sprintf("Invalid key-value pair format: %s", keyPair))
				return false, fmt.Errorf("invalid key-value pair: %s", keyPair)
			}

			key, value := parts[0], parts[1]
			if labelValue, exists := labels[key]; exists && labelValue == value {
				logger.Info(fmt.Sprintf("Cluster %s matches label criteria: %s=%s", clusterName, key, value))
				return true, nil
			}
		}
	}

	logger.Info(fmt.Sprintf("No matching labels found for cluster %s", clusterName))
	return false, nil
}
