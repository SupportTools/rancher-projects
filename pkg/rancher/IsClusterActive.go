package rancher

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func IsClusterActive(clusterName string, cfg *config.Config) bool {
	fmt.Printf("Checking if cluster %s is active...\n", clusterName)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
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
		log.Printf("Cluster %s is not active. Status code: %d\n", clusterName, resp.StatusCode)
		return false
	}

	// Return true if the cluster is active
	return true
}
