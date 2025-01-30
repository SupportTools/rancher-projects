package main

import (
	"github.com/supporttools/rancher-projects/pkg/config"
	"github.com/supporttools/rancher-projects/pkg/logging"
	"github.com/supporttools/rancher-projects/pkg/rancher"
)

var logger = logging.SetupLogging()

func main() {
	// Initialize configuration
	config.Init()

	// Get the configuration
	cfg := config.GetConfig()

	if cfg.ShowHelp {
		config.PrintHelp()
		return
	}

	logger.Info("Starting Rancher-Projects...")

	// Verify access to Rancher
	logger.Info("Verifying access to Rancher...")
	if err := rancher.VerifyAccess(cfg); err != nil {
		logger.Error("Failed to verify access to Rancher: ", err)
		return
	}

	// Determine if handling a single cluster or multiple clusters
	if cfg.ClusterType == "" && cfg.ClusterLabels == "" {
		logger.Info("Processing a single cluster...")
		if err := rancher.SingleCluster(cfg); err != nil {
			logger.Error("Failed to handle single cluster: ", err)
			return
		}
	} else {
		logger.Info("Processing multiple clusters...")
		if err := rancher.MultiCluster(cfg); err != nil {
			logger.Error("Failed to handle multiple clusters: ", err)
			return
		}
	}
}
