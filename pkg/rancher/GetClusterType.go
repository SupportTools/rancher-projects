package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func GetClusterType(cfg *config.Config, clusterName string) string {
	fmt.Println("Getting cluster type...")
	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster type: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster type. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster type: %v", err)
	}

	if len(clusters) > 0 {
		if provider, ok := clusters[0]["provider"].(string); ok {
			return provider
		}
	}

	return ""
}
