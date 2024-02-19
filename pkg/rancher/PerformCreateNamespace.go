package rancher

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func PerformCreateNamespace(cfg *config.Config, clusterID, namespace string) {
	fmt.Println("Checking if namespace", namespace, "exists...")
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces", cfg.RancherServerURL, clusterID)
	reqBody := fmt.Sprintf(`{"type": "namespace", "metadata": {"name": "%s"}}`, namespace)

	req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to check if namespace %s exists: %v", namespace, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Successfully created namespace", namespace)
		fmt.Println("Sleeping for 5 seconds to allow namespace to settle...")
		time.Sleep(5 * time.Second)
	} else if resp.StatusCode == http.StatusConflict {
		fmt.Println("Namespace", namespace, "already exists")
	} else {
		log.Fatalf("Failed to create namespace %s. Status code: %d", namespace, resp.StatusCode)
	}
}
