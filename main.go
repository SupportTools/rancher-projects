package main

import (
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
	rancher.VerifyAccess(cfg)

	if config.GetConfig().ClusterType == "" && config.GetConfig().ClusterLabels == "" {
		rancher.HandleSingleCluster(cfg)
	} else {
		rancher.HandleMultiCluster(cfg)
	}
}
