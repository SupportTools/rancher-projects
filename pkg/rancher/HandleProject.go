package rancher

import "github.com/supporttools/rancher-projects/pkg/config"

func HandleProject(cfg *config.Config, clusterID string) {
	VerifyProject(cfg, clusterID, cfg.ProjectName)

	if cfg.ProjectName != "" {
		VerifyNamespace(cfg, clusterID, cfg.ProjectName)
	}

	if cfg.CreateKubeconfig {
		GenerateKubeconfig(cfg, cfg.KubeconfigFile, clusterID)
	}
}
