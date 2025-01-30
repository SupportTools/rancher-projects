package rancher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyNamespace checks if a given namespace exists within a specified cluster.
func VerifyNamespace(cfg *config.Config, clusterID, namespace string) error {
	logger.Info(fmt.Sprintf("Verifying namespace %s in cluster %s...", namespace, clusterID))

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request for namespace %s: %v", namespace, err))
		return fmt.Errorf("failed to create HTTP request for namespace %s: %v", namespace, err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 10 * time.Second, // More idiomatic timeout setup
	}

	logger.Info(fmt.Sprintf("Sending GET request to verify namespace %s in cluster %s...", namespace, clusterID))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to verify namespace %s: %v", namespace, err))
		return fmt.Errorf("failed to verify namespace %s: %v", namespace, err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		logger.Info(fmt.Sprintf("Successfully found namespace %s in cluster %s.", namespace, clusterID))
		return nil
	case http.StatusNotFound:
		logger.Error(fmt.Sprintf("Namespace %s not found in cluster %s", namespace, clusterID))
		return fmt.Errorf("namespace %s not found in cluster %s", namespace, clusterID)
	default:
		logger.Error(fmt.Sprintf("Unexpected status code %d while verifying namespace %s in cluster %s", resp.StatusCode, namespace, clusterID))
		return fmt.Errorf("unexpected status code %d while verifying namespace %s in cluster %s", resp.StatusCode, namespace, clusterID)
	}
}
