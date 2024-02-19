package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func GetClusterStatus(cfg *config.Config, clusterName string) string {
	fmt.Println("Getting cluster status...")
	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster status. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster status: %v", err)
	}

	if len(clusters) > 0 {
		return clusters[0]["state"].(string)
	}

	return ""
}
