package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetClusterID fetches the cluster ID for a given cluster name from Rancher.
func GetClusterID(cfg *config.Config) (string, error) {
	fmt.Printf("Getting cluster ID for %s...\n", cfg.ClusterName)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, cfg.ClusterName)
	// Updated to use http.NoBody
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get cluster ID, status code: %d", resp.StatusCode)
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode JSON response: %w", err)
	}

	if len(data.Data) == 0 {
		return "", fmt.Errorf("failed to find cluster ID for cluster name: %s", cfg.ClusterName)
	}

	clusterID := data.Data[0].ID
	fmt.Printf("Cluster ID: %s\n", clusterID)
	return clusterID, nil
}
