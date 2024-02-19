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
	fmt.Printf("Verifying project %s...\n", projectName)
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for project %s: %v", projectName, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify project %s: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to verify project %s. Status code: %d", projectName, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body for project %s: %v", projectName, err)
	}

	var projects []map[string]interface{}
	if err = json.Unmarshal(body, &projects); err != nil {
		return fmt.Errorf("failed to parse project data for %s: %v", projectName, err)
	}

	if len(projects) == 0 {
		return fmt.Errorf("failed to find project %s", projectName)
	}

	fmt.Printf("Successfully found project %s\n", projectName)
	return nil
}
