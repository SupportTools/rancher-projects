package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetClusterStatus fetches the status of a specified cluster from Rancher.
func GetClusterStatus(cfg *config.Config, clusterName string) (string, error) {
	fmt.Println("Getting cluster status...")
	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)

	// Use http.NoBody instead of nil
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve cluster status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve cluster status. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var clusters []map[string]interface{}
	if err = json.Unmarshal(body, &clusters); err != nil {
		return "", fmt.Errorf("failed to parse cluster status: %v", err)
	}

	if len(clusters) == 0 {
		return "", fmt.Errorf("no clusters found with name: %s", clusterName)
	}

	status, ok := clusters[0]["state"].(string)
	if !ok {
		return "", fmt.Errorf("cluster status not found or is not a string")
	}

	return status, nil
}
