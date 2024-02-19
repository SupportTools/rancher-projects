package rancher

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func VerifyCluster(cfg *config.Config) {
	fmt.Printf("Verifying cluster %s...\n", cfg.ClusterName)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, cfg.ClusterName)
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
		fmt.Printf("Failed to find cluster %s\n", cfg.ClusterName)
		os.Exit(2)
	}

	fmt.Printf("Successfully found cluster %s\n", cfg.ClusterName)
}
