package rancher

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// AssignNamespaceToProject updates the project ID associated with a namespace in Rancher.
// It returns an error in case of failure else nil.
func AssignNamespaceToProject(cfg *config.Config, clusterID, namespace, projectID string) error {
	logger.Info(fmt.Sprintf("Assigning namespace %s to project %s in cluster %s...", namespace, cfg.ProjectName, clusterID))

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	logger.Debug(fmt.Sprintf("Generated URL for namespace update: %s", url))

	// Use http.NoBody instead of nil for requests with no body.
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP GET request: %v", err))
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Debug("Sending GET request to fetch namespace details...")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("HTTP request failed: %v", err))
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Unexpected response status: %d", resp.StatusCode))
		return fmt.Errorf("failed to assign namespace to project, status code: %d", resp.StatusCode)
	}

	logger.Debug("Decoding JSON response from Rancher API...")
	var namespaceData struct {
		Metadata struct {
			Annotations struct {
				ProjectID string `json:"field.cattle.io/projectId"`
			} `json:"annotations"`
		} `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&namespaceData)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to decode JSON response: %v", err))
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	logger.Info(fmt.Sprintf("Current project ID for namespace %s: %s", namespace, namespaceData.Metadata.Annotations.ProjectID))
	namespaceData.Metadata.Annotations.ProjectID = projectID

	payload, err := json.Marshal(namespaceData)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal JSON payload: %v", err))
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	logger.Debug(fmt.Sprintf("Generated payload for namespace update: %s", string(payload)))
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP PUT request: %v", err))
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	logger.Info(fmt.Sprintf("Sending PUT request to update namespace %s to project %s...", namespace, cfg.ProjectName))
	resp, err = client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP PUT request: %v", err))
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("Failed to assign namespace to project, status code: %d", resp.StatusCode))
		return fmt.Errorf("failed to assign namespace to project, status code: %d", resp.StatusCode)
	}

	logger.Info(fmt.Sprintf("Successfully assigned namespace %s to project %s", namespace, cfg.ProjectName))
	return nil
}
