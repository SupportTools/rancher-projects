package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// CreateProject checks for the existence of a project by name within a cluster and creates it if not found.
func CreateProject(cfg *config.Config, clusterID, projectName string) error {
	fmt.Printf("Checking if project %s exists...\n", projectName)
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to check if project %s exists: %v", projectName, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check if project %s exists: %v", projectName, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("Project %s already exists\n", projectName)
		return nil
	case http.StatusNotFound:
		fmt.Printf("Creating project %s...\n", projectName)
		createProjectURL := fmt.Sprintf("%s/v3/projects", cfg.RancherServerURL)
		projectData := map[string]string{
			"type":      "project",
			"name":      projectName,
			"clusterId": clusterID,
		}
		reqBodyBytes, _ := json.Marshal(projectData)

		req, err = http.NewRequest("POST", createProjectURL, bytes.NewReader(reqBodyBytes))
		if err != nil {
			return fmt.Errorf("failed to create request for new project %s: %v", projectName, err)
		}
		req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to create project %s: %v", projectName, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to create project %s. Status code: %d", projectName, resp.StatusCode)
		}

		fmt.Printf("Successfully created project %s\n", projectName)
		return nil
	default:
		return fmt.Errorf("unexpected status code %d while checking project %s", resp.StatusCode, projectName)
	}
}
