package rancher

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyAccess checks if the provided credentials have access to the Rancher server.
func VerifyAccess(cfg *config.Config) error {
	fmt.Println("Verifying access to Rancher server...")

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	url := fmt.Sprintf("%s/v3/", cfg.RancherServerURL)

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to authenticate to %s with status code: %d", cfg.RancherServerURL, resp.StatusCode)
	}

	fmt.Printf("Successfully authenticated to %s\n", cfg.RancherServerURL)
	return nil
}
