package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	ClusterName           string
	ClusterType           string
	ClusterLabels         string
	ClusterStatus         string
	ClusterID             string
	ClusterIDs            []string
	CreateKubeconfig      bool
	CreateNamespace       bool
	CreateProject         bool
	FilterClustersByType  bool
	FilterClustersByLabel bool
	KubeconfigFile        string
	KubeconfigDir         string
	KubeconfigPrefix      string
	Namespace             string
	ProjectName           string
	RancherAccessKey      string
	RancherSecretKey      string
	RancherServerURL      string
	Debug                 bool
	ShowHelp              bool
}

var cfg Config

var (
	// Define the currentConfig instance
	currentConfig = &Config{}
)

func Init() {
	flag.BoolVar(&currentConfig.ShowHelp, "h", false, "Show help message")
	flag.StringVar(&currentConfig.ClusterName, "cluster-name", "", "The name of the cluster")
	flag.BoolVar(&currentConfig.CreateKubeconfig, "create-kubeconfig", false, "Generate Kubeconfig")
	flag.BoolVar(&currentConfig.CreateNamespace, "create-namespace", false, "Create a namespace")
	flag.BoolVar(&currentConfig.CreateProject, "create-project", false, "Create a project")
	flag.BoolVar(&currentConfig.FilterClustersByType, "get-clusters-by-type", false, "Get clusters by type")
	flag.BoolVar(&currentConfig.FilterClustersByLabel, "get-clusters-by-label", false, "Get clusters by label")
	flag.StringVar(&currentConfig.KubeconfigFile, "kubeconfig", "", "Kubeconfig file")
	flag.StringVar(&currentConfig.KubeconfigDir, "kubeconfig-dir", "", "Kubeconfig directory")
	flag.StringVar(&currentConfig.Namespace, "namespace", "", "Namespace")
	flag.StringVar(&currentConfig.ProjectName, "project-name", "", "Project name")
	flag.StringVar(&currentConfig.RancherAccessKey, "rancher-access-key", "", "Rancher access key")
	flag.StringVar(&currentConfig.RancherSecretKey, "rancher-secret-key", "", "Rancher secret key")
	flag.StringVar(&currentConfig.RancherServerURL, "rancher-server", "", "Rancher server URL")
	flag.Parse()

	// Initialize the cfg variable with command line flag values
	cfg = *currentConfig

	// Load additional configuration from environment variables
	LoadConfig()

	// Check for missing required settings
	checkMissingSettings()
}

func checkMissingSettings() {
	requiredFlags := []string{
		"rancher-server",
		"rancher-access-key",
		"rancher-secret-key",
	}

	requiredFlagCombos := [][]string{
		{"get-clusters-by-label", "get-clusters-by-type"},
	}

	missingRequiredFlags := []string{}
	missingRequiredFlagCombos := [][]string{}

	flagSet := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) {
		flagSet[f.Name] = true
	})

	for _, flagName := range requiredFlags {
		if !flagSet[flagName] {
			missingRequiredFlags = append(missingRequiredFlags, "--"+flagName)
		}
	}

	for _, flagCombo := range requiredFlagCombos {
		if (!flagSet[flagCombo[0]] && !flagSet[flagCombo[1]]) ||
			(flagSet[flagCombo[0]] && flagSet[flagCombo[1]]) {
			missingRequiredFlagCombos = append(missingRequiredFlagCombos, flagCombo)
		}
	}

	if len(missingRequiredFlags) > 0 || len(missingRequiredFlagCombos) > 0 {
		fmt.Println("Missing required flags:")
		if len(missingRequiredFlags) > 0 {
			fmt.Println("\nSingle flags:")
			for _, flagName := range missingRequiredFlags {
				fmt.Println("-", flagName)
			}
		}
		if len(missingRequiredFlagCombos) > 0 {
			fmt.Println("\nFlag combinations:")
			for _, flagCombo := range missingRequiredFlagCombos {
				fmt.Printf("- Either %s and %s\n", "--"+flagCombo[0], "--"+flagCombo[1])
			}
		}
		fmt.Println("\nPlease provide the missing flags.")
		PrintHelp()
		os.Exit(1)
	}
}

// LoadConfig loads configuration from environment variables and command line flags
func LoadConfig() {
	cfg.ClusterType = getEnv("CLUSTER_TYPE")
	cfg.ClusterLabels = getEnv("CLUSTER_LABELS")
	cfg.ClusterStatus = getEnv("CLUSTER_STATUS")
	cfg.ClusterID = getEnv("CLUSTER_ID")
	cfg.ClusterIDs = getEnvArray("CLUSTER_IDS", ",")
	cfg.ProjectName = getEnv("PROJECT_NAME")
	cfg.RancherServerURL = getEnv("RANCHER_SERVER")
	cfg.RancherAccessKey = getEnv("RANCHER_ACCESS_KEY")
	cfg.RancherSecretKey = getEnv("RANCHER_SECRET_KEY")
	cfg.KubeconfigFile = getEnv("KUBECONFIG")
	cfg.KubeconfigDir = getEnv("KUBECONFIG_DIR")
	cfg.KubeconfigPrefix = getEnv("KUBECONFIG_PREFIX")
	cfg.Namespace = getEnv("NAMESPACE")
	cfg.Debug = getEnvBool("DEBUG", false)
}

// GetConfig returns the current configuration instance
func GetConfig() *Config {
	return currentConfig
}

// getEnv gets an environment variable.
// If the variable is not found, it returns an empty string.
func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		if flag.Lookup(key) != nil {
			return flag.Lookup(key).DefValue
		}
		log.Fatalf("Environment variable %s is not set\n", key)
	}
	return value
}

// getEnvBool gets an environment variable and returns it as a boolean
func getEnvBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value == "true" || value == "1"
}

// getEnvArray gets an environment variable and returns it as an array
func getEnvArray(key, separator string) []string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return []string{}
	}
	return strings.Split(value, separator)
}

func PrintHelp() {
	fmt.Println("Usage: rancher-projects [options]")
	fmt.Println("Options:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  --%s %s\n", f.Name, f.Usage)
	})

	// Show usage examples
	fmt.Println("\nUsage examples:")
	fmt.Println("  Getting a kubeconfig for a single cluster, create a project and namespace at the same time:")
	fmt.Println("    rancher-projects \\")
	fmt.Println("    --rancher-server \"https://rancher.mattox.local\" \\")
	fmt.Println("    --rancher-access-key \"token-abcde\" \\")
	fmt.Println("    --rancher-secret-key \"123456789abcdefghijklmnopqrstuvwxyz\" \\")
	fmt.Println("    --cluster-name \"MyCluster\" \\")
	fmt.Println("    --project-name \"MyProject\" \\")
	fmt.Println("    --namespace \"mynamespace\" \\")
	fmt.Println("    --create-project true \\")
	fmt.Println("    --create-namespace true \\")
	fmt.Println("    --create-kubeconfig true \\")
	fmt.Println("    --kubeconfig \"rancher-projects-kubeconfig\"")
	fmt.Println("\n  Getting a kubeconfig for multiple RKE2 clusters:")
	fmt.Println("    rancher-projects \\")
	fmt.Println("    --rancher-server \"https://rancher.mattox.local\" \\")
	fmt.Println("    --rancher-access-key \"token-abcde\" \\")
	fmt.Println("    --rancher-secret-key \"123456789abcdefghijklmnopqrstuvwxyz\" \\")
	fmt.Println("    --create-kubeconfig true \\")
	fmt.Println("    --kubeconfig-dir \"~/.kube/\"")
	fmt.Println("    --get-clusters-by-type \"rke2\"")

}

// PrintConfig prints the loaded configuration
func PrintConfig() {
	fmt.Println("Loaded Configuration:")
	fmt.Println("Show Help:", cfg.ShowHelp)
	fmt.Println("Cluster Name:", cfg.ClusterName)
	fmt.Println("Create Kubeconfig:", cfg.CreateKubeconfig)
	fmt.Println("Create Namespace:", cfg.CreateNamespace)
	fmt.Println("Create Project:", cfg.CreateProject)
	fmt.Println("Filter Clusters by Type:", cfg.FilterClustersByType)
	fmt.Println("Filter Clusters by Label:", cfg.FilterClustersByLabel)
	fmt.Println("Kubeconfig File:", cfg.KubeconfigFile)
	fmt.Println("Kubeconfig Dir:", cfg.KubeconfigDir)
	fmt.Println("Kubeconfig Prefix:", cfg.KubeconfigPrefix)
	fmt.Println("Rancher Server URL:", cfg.RancherServerURL)
	fmt.Println("Rancher Access Key:", cfg.RancherAccessKey)
}

// ParseFlags parses command line flags
func ParseFlags() {
	flag.Parse()
}

func (c *Config) GetClusterType() string {
	return c.ClusterType
}

func (c *Config) SetClusterType(clusterType string) {
	c.ClusterType = clusterType
}

func (c *Config) GetClusterLabels() string {
	return c.ClusterLabels
}

func (c *Config) SetClusterLabels(clusterLabels string) {
	c.ClusterLabels = clusterLabels
}

func (c *Config) GetClusterStatus() string {
	return c.ClusterStatus
}

func (c *Config) SetClusterStatus(clusterStatus string) {
	c.ClusterStatus = clusterStatus
}

func (c *Config) GetClusterID() string {
	return c.ClusterID
}

func (c *Config) SetClusterID(clusterID string) {
	c.ClusterID = clusterID
}
