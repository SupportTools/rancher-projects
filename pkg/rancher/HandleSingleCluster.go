package rancher

import "github.com/supporttools/rancher-projects/pkg/config"

func HandleSingleCluster(cfg *config.Config) {
	VerifyCluster(cfg)
	clusterID := GetClusterID(cfg)

	if cfg.ProjectName != "" {
		HandleProject(cfg, clusterID)
	}

	if cfg.CreateKubeconfig {
		GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID)
	}
}
