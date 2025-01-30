package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetProjectInfo fetches the project ID for a given project name within a specified cluster.
func GetProjectInfo(cfg *config.Config, clusterID, projectName string) (string, error) {
	logger.Info(fmt.Sprintf("Fetching project info for project: %s in cluster: %s", projectName, clusterID))

	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	// Updated to use http.NoBody instead of nil
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info("Sending GET request to retrieve project info...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP request: %v", err))
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to get project info, status code: %d", resp.StatusCode))
		return "", fmt.Errorf("failed to get project info, status code: %d", resp.StatusCode)
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode JSON response: %v", err))
		return "", fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if len(data.Data) == 0 {
		logger.Error(fmt.Sprintf("Failed to find project info for project name: %s", projectName))
		return "", fmt.Errorf("failed to find project info for project name: %s", projectName)
	}

	projectID := data.Data[0].ID
	logger.Info(fmt.Sprintf("Successfully retrieved project ID: %s", projectID))
	return projectID, nil
}
