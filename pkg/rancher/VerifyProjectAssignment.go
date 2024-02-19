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
	fmt.Println("Verifying project assignment...")

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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
		return fmt.Errorf("failed to decode JSON response: %v", err)
	}

	if namespaceData.Metadata.Annotations.ProjectID != projectID {
		return fmt.Errorf("project ID mismatch: expected %s, got %s", projectID, namespaceData.Metadata.Annotations.ProjectID)
	}

	fmt.Println("Successfully verified project assignment")
	return nil
}
