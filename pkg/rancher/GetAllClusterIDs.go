package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func GetAllClusterIDs(cfg *config.Config) {
	fmt.Println("Getting all cluster IDs...")
	url := fmt.Sprintf("%s/v3/clusters", cfg.RancherServerURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster IDs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster IDs. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster IDs: %v", err)
	}

	cfg.ClusterIDs = make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterName := cluster["name"].(string)
		clusterID := cluster["id"].(string)
		cfg.ClusterIDs[i] = fmt.Sprintf("%s:%s", clusterName, clusterID)
	}

	fmt.Println("Cluster IDs:", cfg.ClusterIDs)
}
