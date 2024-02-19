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
	fmt.Printf("Assigning namespace %s to project %s...\n", namespace, cfg.ProjectName)

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	// Use http.NoBody instead of nil for requests with no body.
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to assign namespace to project, status code: %d", resp.StatusCode)
	}

	var namespaceData struct {
		Metadata struct {
			Annotations struct {
				ProjectID string `json:"field.cattle.io/projectId"`
			} `json:"annotations"`
		} `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&namespaceData)
	if err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	namespaceData.Metadata.Annotations.ProjectID = projectID

	payload, err := json.Marshal(namespaceData)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON payload: %w", err)
	}

	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to assign namespace to project, status code: %d", resp.StatusCode)
	}

	fmt.Printf("Successfully assigned namespace %s to project %s\n", namespace, cfg.ProjectName)
	return nil
}
