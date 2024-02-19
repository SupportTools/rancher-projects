package rancher

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func VerifyAccess(cfg *config.Config) {
	fmt.Println("Verifying access to Rancher server...")

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	url := fmt.Sprintf("%s/v3/", cfg.RancherServerURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to authenticate to %s\n", cfg.RancherServerURL)
		os.Exit(2)
	}

	fmt.Printf("Successfully authenticated to %s\n", cfg.RancherServerURL)
}
