package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func VerifyProject(cfg *config.Config, clusterID, projectName string) {
	fmt.Println("Verifying project", projectName, "...")
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", cfg.RancherServerURL, clusterID, projectName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to verify project %s: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to verify project %s. Status code: %d", projectName, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var projects []map[string]interface{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		log.Fatalf("Failed to parse project data: %v", err)
	}

	if len(projects) == 0 {
		log.Fatalf("Failed to find project %s", projectName)
	}

	fmt.Println("Successfully found project", projectName)
}
