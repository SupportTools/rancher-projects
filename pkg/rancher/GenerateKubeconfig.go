package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GenerateKubeconfig creates a kubeconfig file for a specified cluster.
func GenerateKubeconfig(cfg *config.Config, kubeconfigFile, clusterID string) error {
	fmt.Println("Generating kubeconfig...")

	// Construct the request URL for generating kubeconfig.
	url := fmt.Sprintf("%s/v3/clusters/%s?action=generateKubeconfig", cfg.RancherServerURL, clusterID)
	req, err := http.NewRequest("POST", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Encode the authentication credentials and set the request headers.
	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	// Send the request using the HTTP client.
	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to generate kubeconfig, status code: %d", resp.StatusCode)
	}

	// Decode the response body to extract the kubeconfig data.
	var data struct {
		Config string `json:"config"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}

	// Create the directory for the kubeconfig file if it doesn't exist.
	if err := os.MkdirAll(filepath.Dir(kubeconfigFile), 0o755); err != nil {
		return fmt.Errorf("failed to create kubeconfig directory: %w", err)
	}

	// Create the kubeconfig file if it doesn't exist.
	if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
		if _, err := os.Create(kubeconfigFile); err != nil {
			return fmt.Errorf("failed to create kubeconfig file: %w", err)
		}
	}

	// Write the kubeconfig data to the specified file.
	if err := os.WriteFile(kubeconfigFile, []byte(data.Config), 0o644); err != nil {
		return fmt.Errorf("failed to write kubeconfig file: %w", err)
	}

	fmt.Printf("Kubeconfig file generated: %s\n", kubeconfigFile)
	return nil
}
