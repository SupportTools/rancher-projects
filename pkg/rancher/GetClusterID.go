package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func GetClusterID(cfg *config.Config) string {
	fmt.Printf("Getting cluster ID for %s...\n", cfg.ClusterName)

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
		log.Fatal("Failed to get cluster ID")
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	if len(data.Data) == 0 {
		log.Fatal("Failed to find cluster ID")
	}

	clusterID := data.Data[0].ID
	fmt.Printf("Cluster ID: %s\n", clusterID)
	fmt.Println("Successfully got cluster ID")
	return clusterID
}
