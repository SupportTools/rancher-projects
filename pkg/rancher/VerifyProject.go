package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyProject checks if a given project exists within a specified cluster.
func VerifyProject(cfg *config.Config, clusterID, projectName string) error {
	logger.Info(fmt.Sprintf("Verifying project %s in cluster %s...", projectName, clusterID))

	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request for project %s: %v", projectName, err))
		return fmt.Errorf("failed to create HTTP request for project %s: %v", projectName, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info(fmt.Sprintf("Sending GET request to verify project %s in cluster %s...", projectName, clusterID))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to verify project %s: %v", projectName, err))
		return fmt.Errorf("failed to verify project %s: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to verify project %s. Status code: %d", projectName, resp.StatusCode))
		return fmt.Errorf("failed to verify project %s. Status code: %d", projectName, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to read response body for project %s: %v", projectName, err))
		return fmt.Errorf("failed to read response body for project %s: %v", projectName, err)
	}
	logger.Debug(fmt.Sprintf("Received response body: %s", string(body)))

	var response struct { // Define a more structured response if possible
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}
	if err = json.Unmarshal(body, &response); err != nil {
		logger.Error(fmt.Sprintf("Failed to parse project data for %s: %v", projectName, err))
		return fmt.Errorf("failed to parse project data for %s: %v", projectName, err)
	}

	if len(response.Data) == 0 {
		logger.Error(fmt.Sprintf("Project %s not found in cluster %s", projectName, clusterID))
		return fmt.Errorf("project %s not found in cluster %s", projectName, clusterID)
	}

	logger.Info(fmt.Sprintf("Successfully verified project %s exists in cluster %s.", projectName, clusterID))
	return nil
}
