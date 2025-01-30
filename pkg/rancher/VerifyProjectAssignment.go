package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyProjectAssignment checks if a namespace is assigned to the specified project.
func VerifyProjectAssignment(cfg *config.Config, clusterID, namespace, projectID string) error {
	logger.Info(fmt.Sprintf("Verifying project assignment for namespace %s in cluster %s...", namespace, clusterID))

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info(fmt.Sprintf("Sending GET request to verify project assignment for namespace %s...", namespace))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP request: %v", err))
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to verify project assignment, status code: %d", resp.StatusCode))
		return fmt.Errorf("failed to verify project assignment, status code: %d", resp.StatusCode)
	}

	var namespaceData struct {
		Metadata struct {
			Annotations struct {
				ProjectID string `json:"field.cattle.io/projectId"`
			} `json:"annotations"`
		} `json:"metadata"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&namespaceData); err != nil {
		logger.Error(fmt.Sprintf("Failed to decode JSON response: %v", err))
		return fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if namespaceData.Metadata.Annotations.ProjectID != projectID {
		logger.Error(fmt.Sprintf("Project ID mismatch: expected %s, got %s", projectID, namespaceData.Metadata.Annotations.ProjectID))
		return fmt.Errorf("project ID mismatch: expected %s, got %s", projectID, namespaceData.Metadata.Annotations.ProjectID)
	}

	logger.Info(fmt.Sprintf("Successfully verified project assignment for namespace %s in cluster %s.", namespace, clusterID))
	return nil
}
