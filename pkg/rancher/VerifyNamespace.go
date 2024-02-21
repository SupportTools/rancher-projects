package rancher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyNamespace checks if a given namespace exists within a specified cluster.
func VerifyNamespace(cfg *config.Config, clusterID, namespace string) error {
	fmt.Printf("Verifying namespace %s...\n", namespace)
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request for namespace %s: %v", namespace, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second, // More idiomatic timeout setup
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify namespace %s: %v", namespace, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		fmt.Printf("Successfully found namespace %s in cluster %s.\n", namespace, clusterID)
		return nil
	case http.StatusNotFound:
		return fmt.Errorf("namespace %s not found in cluster %s", namespace, clusterID)
	default:
		return fmt.Errorf("unexpected status code %d while verifying namespace %s in cluster %s", resp.StatusCode, namespace, clusterID)
	}
}
