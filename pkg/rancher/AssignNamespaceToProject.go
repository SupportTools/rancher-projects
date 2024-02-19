package rancher

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func AssignNamespaceToProject(cfg *config.Config, clusterID, namespace, projectID string) {
	fmt.Printf("Assigning namespace %s to project %s...\n", namespace, cfg.ProjectName)

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
		log.Fatal("Failed to assign namespace to project")
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

	namespaceData.Metadata.Annotations.ProjectID = projectID

	payload, err := json.Marshal(namespaceData)
	if err != nil {
		log.Fatal("Failed to marshal JSON payload:", err)
	}

	url = fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to assign namespace to project")
	}

	fmt.Printf("Successfully assigned namespace %s to project %s\n", namespace, cfg.ProjectName)
}
