package rancher

import (
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func VerifyNamespace(cfg *config.Config, clusterID, namespace string) {
	fmt.Println("Verifying namespace", namespace, "...")
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", cfg.RancherServerURL, clusterID, namespace)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to verify namespace %s: %v", namespace, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Fatalf("Failed to find namespace %s", namespace)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to verify namespace %s. Status code: %d", namespace, resp.StatusCode)
	}

	fmt.Println("Successfully found namespace", namespace)
}
