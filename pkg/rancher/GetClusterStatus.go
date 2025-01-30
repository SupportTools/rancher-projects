package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetClusterStatus fetches the status of a specified cluster from Rancher.
func GetClusterStatus(cfg *config.Config, clusterName string) (string, error) {
	logger.Info(fmt.Sprintf("Fetching cluster status for cluster: %s", clusterName))

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	// Use http.NoBody instead of nil
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info("Sending GET request to retrieve cluster status...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to retrieve cluster status: %v", err))
		return "", fmt.Errorf("failed to retrieve cluster status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to retrieve cluster status. Status code: %d", resp.StatusCode))
		return "", fmt.Errorf("failed to retrieve cluster status. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body: %v", err))
		return "", fmt.Errorf("failed to read response body: %v", err)
	}
	logger.Debug(fmt.Sprintf("Received response body: %s", string(body)))

	var clusters struct {
		Data []struct {
			State string `json:"state"`
		} `json:"data"`
	}

	if err = json.Unmarshal(body, &clusters); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse cluster status: %v", err))
		return "", fmt.Errorf("failed to parse cluster status: %v", err)
	}

	if len(clusters.Data) == 0 {
		logger.Error(fmt.Sprintf("No clusters found with name: %s", clusterName))
		return "", fmt.Errorf("no clusters found with name: %s", clusterName)
	}

	status := clusters.Data[0].State
	logger.Info(fmt.Sprintf("Cluster %s status: %s", clusterName, status))
	return status, nil
}
