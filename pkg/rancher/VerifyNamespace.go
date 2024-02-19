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
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to verify namespace %s: %v", namespace, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("failed to find namespace %s", namespace)
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to verify namespace %s. Status code: %d", namespace, resp.StatusCode)
	}

	fmt.Printf("Successfully found namespace %s\n", namespace)
	return nil
}
