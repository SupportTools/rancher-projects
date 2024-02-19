package rancher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func FilterByClusterLabels(clusterName string, keyPairs []string, cfg *config.Config) bool {
	fmt.Println("Filtering by cluster label")
	fmt.Println("Cluster name:", clusterName)
	fmt.Println("Labels:", keyPairs)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(cfg.RancherAccessKey, cfg.RancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to filter clusters by label: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to filter clusters by label. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse clusters by label: %v", err)
	}

	for _, cluster := range clusters {
		labels := cluster["labels"].(map[string]interface{})
		for _, keyPair := range keyPairs {
			parts := strings.Split(keyPair, "=")
			if len(parts) != 2 {
				log.Fatalf("Invalid key-value pair: %s", keyPair)
			}

			key := parts[0]
			value := parts[1]

			labelValue, exists := labels[key]
			if exists && labelValue == value {
				return true
			}
		}
	}

	return false
}
