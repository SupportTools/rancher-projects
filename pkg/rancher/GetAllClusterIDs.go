package rancher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// GetAllClusterIDs fetches all cluster IDs from Rancher and updates the configuration accordingly.
func GetAllClusterIDs(cfg *config.Config) error {
	fmt.Println("Getting all cluster IDs...")
	url := fmt.Sprintf("%s/v3/clusters", cfg.RancherServerURL)

	// Updated to use http.NoBody
	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to retrieve cluster IDs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve cluster IDs. Status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	var clusters []map[string]interface{}
	if err = json.Unmarshal(body, &clusters); err != nil {
		return fmt.Errorf("failed to parse cluster IDs: %v", err)
	}

	cfg.ClusterIDs = make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterName, okName := cluster["name"].(string)
		clusterID, okID := cluster["id"].(string)
		if !okName || !okID {
			continue // Skip if types do not match expectations
		}
		cfg.ClusterIDs[i] = fmt.Sprintf("%s:%s", clusterName, clusterID)
	}

	fmt.Println("Cluster IDs:", cfg.ClusterIDs)
	return nil
}
