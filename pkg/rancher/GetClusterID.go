package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetClusterID fetches the cluster ID for a given cluster name from Rancher.
func GetClusterID(cfg *config.Config) (string, error) {
	logger.Info(fmt.Sprintf("Fetching cluster ID for cluster: %s", cfg.ClusterName))

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, cfg.ClusterName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	// Updated to use http.NoBody
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info("Sending GET request to retrieve cluster ID...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP request: %v", err))
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to get cluster ID, status code: %d", resp.StatusCode))
		return "", fmt.Errorf("failed to get cluster ID, status code: %d", resp.StatusCode)
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode JSON response: %v", err))
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(data.Data) == 0 {
		logger.Error(fmt.Sprintf("Failed to find cluster ID for cluster name: %s", cfg.ClusterName))
		return "", fmt.Errorf("failed to find cluster ID for cluster name: %s", cfg.ClusterName)
	}

	clusterID := data.Data[0].ID
	logger.Info(fmt.Sprintf("Successfully retrieved cluster ID: %s", clusterID))
	return clusterID, nil
}
