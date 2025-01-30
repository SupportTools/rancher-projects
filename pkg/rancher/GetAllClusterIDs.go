package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetAllClusterIDs fetches all cluster IDs from Rancher and updates the configuration accordingly.
func GetAllClusterIDs(cfg *config.Config) error {
	logger.Info("Fetching all cluster IDs from Rancher...")
	url := fmt.Sprintf("%s/v3/clusters", cfg.RancherServerURL)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	// Updated to use http.NoBody
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info("Sending GET request to fetch cluster IDs...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to retrieve cluster IDs: %v", err))
		return fmt.Errorf("failed to retrieve cluster IDs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to retrieve cluster IDs. Status code: %d", resp.StatusCode))
		return fmt.Errorf("failed to retrieve cluster IDs. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return fmt.Errorf("failed to read response body: %v", err)
	}
	logger.Debug(fmt.Sprintf("Received response body: %s", string(body)))

	var clusters []map[string]interface{}
	if err = json.Unmarshal(body, &clusters); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse cluster IDs: %v", err))
		return fmt.Errorf("failed to parse cluster IDs: %v", err)
	}

	cfg.ClusterIDs = make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterName, okName := cluster["name"].(string)
		clusterID, okID := cluster["id"].(string)
		if !okName || !okID {
			logger.Warn(fmt.Sprintf("Skipping cluster due to missing or invalid name/id: %v", cluster))
			continue // Skip if types do not match expectations
		}
		cfg.ClusterIDs[i] = fmt.Sprintf("%s:%s", clusterName, clusterID)
	}

	logger.Info(fmt.Sprintf("Retrieved Cluster IDs: %v", cfg.ClusterIDs))
	return nil
}
