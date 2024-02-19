package main

import (
	"log"

	"github.com/supporttools/rancher-projects/pkg/config"
	"github.com/supporttools/rancher-projects/pkg/rancher"
)

func main() {
	// Initialize configuration
	config.Init()

	// Get the configuration
	cfg := config.GetConfig()

	if cfg.ShowHelp {
		config.PrintHelp()
		return
	}

	// Verify access to Rancher
	if err := rancher.VerifyAccess(cfg); err != nil {
		log.Fatalf("Failed to verify access to Rancher: %v", err)
	}

	// Depending on the configuration, handle a single cluster or multiple clusters
	if cfg.ClusterType == "" && cfg.ClusterLabels == "" {
		if err := rancher.SingleCluster(cfg); err != nil {
			log.Fatalf("Failed to handle single cluster: %v", err)
		}
	} else {
		if err := rancher.MultiCluster(cfg); err != nil {
			log.Fatalf("Failed to handle multiple clusters: %v", err)
		}
	}
}
