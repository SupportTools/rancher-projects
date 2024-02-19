package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/supporttools/rancher-projects/pkg/config"
)

func GenerateKubeconfig(cfg *config.Config, kubeconfigFile, clusterID string) {
	fmt.Println("Generating kubeconfig...")

	url := fmt.Sprintf("%s/v3/clusters/%s?action=generateKubeconfig", cfg.RancherServerURL, clusterID)
	req, err := http.NewRequest("POST", url, nil)
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
		log.Fatal("Failed to generate kubeconfig")
	}

	var data struct {
		Config string `json:"config"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	err = ioutil.WriteFile(kubeconfigFile, []byte(data.Config), 0644)
	if err != nil {
		log.Fatal("Failed to write kubeconfig file:", err)
	}

	fmt.Printf("Kubeconfig file generated: %s\n", kubeconfigFile)
	fmt.Println("Successfully generated kubeconfig")
}
