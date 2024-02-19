package rancher

import (
    "encoding/base64"
    "fmt"
    "net/http"
    "time"

    "github.com/supporttools/rancher-projects/pkg/config"
)

// IsClusterActive checks if a specified cluster is active within Rancher.
func IsClusterActive(clusterName string, cfg *config.Config) (bool, error) {
    fmt.Printf("Checking if cluster %s is active...\n", clusterName)

    url := fmt.Sprintf("%s/v3/clusters?name=%s", cfg.RancherServerURL, clusterName)
    // Updated to use http.NoBody instead of nil
    req, err := http.NewRequest("GET", url, http.NoBody)
    if err != nil {
        return false, fmt.Errorf("failed to create HTTP request: %v", err)
    }

    authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", cfg.RancherAccessKey, cfg.RancherSecretKey)))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

    client := &http.Client{
        Timeout: time.Second * 10, // 10 seconds timeout
    }
    resp, err := client.Do(req)
    if err != nil {
        return false, fmt.Errorf("failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Adjusted to provide a more general error message, as non-OK status doesn't necessarily mean inactive
        return false, fmt.Errorf("failed to check cluster status for '%s', status code: %d", clusterName, resp.StatusCode)
    }

    // Return true if the cluster is active, with no error
    return true, nil
}
