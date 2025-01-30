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
	logger.Info(fmt.Sprintf("Starting CreateProject for project %s in cluster %s", projectName, clusterID))

	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)
	logger.Debug(fmt.Sprintf("Generated request URL for project check: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create GET request for project %s: %v", projectName, err))
		return fmt.Errorf("failed to check if project %s exists: %v", projectName, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info(fmt.Sprintf("Sending GET request to check if project %s exists...", projectName))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to execute GET request for project %s: %v", projectName, err))
		return fmt.Errorf("failed to check if project %s exists: %v", projectName, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info(fmt.Sprintf("Project %s already exists", projectName))
		return nil
	case http.StatusNotFound:
		logger.Info(fmt.Sprintf("Project %s not found, proceeding to create it...", projectName))

		createProjectURL := fmt.Sprintf("%s/v3/projects", cfg.RancherServerURL)
		projectData := map[string]string{
			"type":      "project",
			"name":      projectName,
			"clusterId": clusterID,
		}

		reqBodyBytes, err := json.Marshal(projectData)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to marshal project data for %s: %v", projectName, err))
			return fmt.Errorf("failed to marshal project data for %s: %v", projectName, err)
		}
		logger.Debug(fmt.Sprintf("Generated request body for project creation: %s", string(reqBodyBytes)))

		req, err = http.NewRequest("POST", createProjectURL, bytes.NewReader(reqBodyBytes))
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to create POST request for new project %s: %v", projectName, err))
			return fmt.Errorf("failed to create request for new project %s: %v", projectName, err)
		}
		req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
		req.Header.Set("Content-Type", "application/json")

		logger.Info(fmt.Sprintf("Sending POST request to create project %s...", projectName))
		resp, err = client.Do(req)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to execute POST request for project %s: %v", projectName, err))
			return fmt.Errorf("failed to create project %s: %v", projectName, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			logger.Error(fmt.Sprintf("Failed to create project %s. Status code: %d", projectName, resp.StatusCode))
			return fmt.Errorf("failed to create project %s. Status code: %d", projectName, resp.StatusCode)
		}

		logger.Info(fmt.Sprintf("Successfully created project %s", projectName))
		return nil
	default:
		logger.Error(fmt.Sprintf("Unexpected status code %d while checking project %s", resp.StatusCode, projectName))
		return fmt.Errorf("unexpected status code %d while checking project %s", resp.StatusCode, projectName)
	}
}
