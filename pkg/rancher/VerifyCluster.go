package rancher

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/supporttools/rancher-projects/pkg/config"
)

// VerifyCluster checks if a specified cluster exists within Rancher.
func VerifyCluster(cfg *config.Config) error {
	logger.Info(fmt.Sprintf("Verifying cluster %s...", cfg.ClusterName))

	url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, cfg.ClusterName)
	logger.Debug(fmt.Sprintf("Generated request URL: %s", url))

	req, err := http.NewRequest("GET", url, http.NoBody)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to create HTTP request: %v", err))
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{
		Timeout: time.Second * 10, // 10 seconds timeout
	}

	logger.Info(fmt.Sprintf("Sending GET request to verify cluster %s...", cfg.ClusterName))
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to send HTTP request: %v", err))
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var responseBody struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&responseBody); err == nil {
			logger.Error(fmt.Sprintf("Failed to find cluster %s: %s", cfg.ClusterName, responseBody.Message))
			return fmt.Errorf("failed to find cluster %s: %s", cfg.ClusterName, responseBody.Message)
		}
		logger.Error(fmt.Sprintf("Failed to find cluster %s, status code: %d", cfg.ClusterName, resp.StatusCode))
		return fmt.Errorf("failed to find cluster %s, status code: %d", cfg.ClusterName, resp.StatusCode)
	}

	logger.Info(fmt.Sprintf("Successfully found cluster %s", cfg.ClusterName))
	return nil
}
