package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func VerifyProjectAssignment(cfg *config.Config, clusterID, namespace, projectID string) {
	fmt.Println("Verifying project assignment...")

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to verify project assignment")
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
		log.Fatal("Failed to decode JSON response:", err)
	}

	if namespaceData.Metadata.Annotations.ProjectID != projectID {
		log.Fatal("Failed to verify project assignment")
	}

	fmt.Println("Successfully verified project assignment")
}
