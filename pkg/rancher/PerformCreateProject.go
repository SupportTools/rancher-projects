package rancher

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func PerformCreateProject(cfg *config.Config, clusterID, projectName string) {
	fmt.Println("Checking if project", projectName, "exists...")
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to check if project %s exists: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Project", projectName, "already exists")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Creating project", projectName, "...")
		createProjectURL := fmt.Sprintf("%s/v3/projects", cfg.RancherServerURL)
		reqBody := fmt.Sprintf(`{"type": "project", "name": "%s", "clusterId": "%s"}`, projectName, clusterID)

		req, _ = http.NewRequest("POST", createProjectURL, strings.NewReader(reqBody))
		req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("Failed to create project %s: %v", projectName, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf("Failed to create project %s. Status code: %d", projectName, resp.StatusCode)
		}

		fmt.Println("Successfully created project", projectName)
	} else {
		log.Fatalf("Failed to check project %s. Status code: %d", projectName, resp.StatusCode)
	}
}
