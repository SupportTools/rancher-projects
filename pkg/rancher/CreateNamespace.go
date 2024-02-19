package rancher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// CreateNamespace attempts to create a namespace within a specified cluster.
// It waits for 5 seconds if the namespace is successfully created to allow it to settle.
func CreateNamespace(cfg *config.Config, clusterID, namespace string) error {
	fmt.Printf("Checking if namespace %s exists...\n", namespace)
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces", cfg.RancherServerURL, clusterID)

	namespaceData := map[string]interface{}{
		"type": "namespace",
		"metadata": map[string]string{
			"name": namespace,
		},
	}
	reqBodyBytes, err := json.Marshal(namespaceData)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for namespace %s: %v", namespace, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to check if namespace %s exists: %v", namespace, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusCreated:
		fmt.Printf("Successfully created namespace %s\n", namespace)
		fmt.Println("Sleeping for 5 seconds to allow namespace to settle...")
		time.Sleep(5 * time.Second)
		return nil
	case http.StatusConflict:
		fmt.Printf("Namespace %s already exists\n", namespace)
		return nil
	default:
		return fmt.Errorf("failed to create namespace %s. Status code: %d", namespace, resp.StatusCode)
	}
}
