package rancher

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// IsClusterActive checks if a specified cluster is active within Rancher.
func IsClusterActive(clusterName string, cfg *config.Config) (bool, error) {
	logger.Info(fmt.Sprintf("Checking if cluster %s is active...", clusterName))

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	// Updated to use http.NoBody instead of nil
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return false, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info(fmt.Sprintf("Sending GET request to check status of cluster %s...", clusterName))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP request: %v", err))
		return false, fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to check cluster status for '%s', status code: %d", clusterName, resp.StatusCode))
		return false, fmt.Errorf("failed to check cluster status for '%s', status code: %d", clusterName, resp.StatusCode)
	}

	logger.Info(fmt.Sprintf("Cluster %s is active.", clusterName))
	return true, nil
}
