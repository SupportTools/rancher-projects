package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	clusterName        string
	clusterType        string
	clusterLabels      string
	clusterStatus      string
	clusterID          string
	clusterIDs         []string
	projectName        string
	createProject      bool
	namespace          string
	createNamespace    bool
	rancherServer      string
	rancherAccessKey   string
	rancherSecretKey   string
	createKubeconfig   bool
	getClustersByType  string
	getClustersByLabel string
	kubeconfig         string
	kubeconfigFile     string
	kubeconfigDir      string
	kubeconfigPrefix   string
	showHelp           bool
	debug              bool
)

type Config struct {
	ClusterName           string
	CreateKubeconfig      bool
	CreateNamespace       bool
	CreateProject         bool
	FilterClustersByType  string
	FilterClustersByLabel string
	KubeconfigFile        string
	Namespace             string
	ProjectName           string
	RancherAccessKey      string
	RancherSecretKey      string
	RancherServerURL      string
}

var config Config

func init() {
	flag.BoolVar(&createKubeconfig, "create-kubeconfig", false, "Create Kubeconfig")
	flag.BoolVar(&createNamespace, "create-namespace", false, "Create Namespace")
	flag.BoolVar(&createProject, "create-project", false, "Create Project")
	flag.StringVar(&clusterType, "get-clusters-by-type", "", "Filter Clusters by Type")
	flag.StringVar(&clusterLabels, "get-clusters-by-label", "", "Filter Clusters by Label")
	flag.StringVar(&kubeconfigFile, "kubeconfig", "rancher-projects-kubeconfig", "Kubeconfig File")
	flag.StringVar(&namespace, "namespace", "", "Namespace")
	flag.StringVar(&projectName, "project-name", "", "Project Name")
	flag.StringVar(&rancherAccessKey, "rancher-access-key", "", "Rancher Access Key")
	flag.StringVar(&rancherSecretKey, "rancher-secret-key", "", "Rancher Secret Key")
	flag.StringVar(&rancherServer, "rancher-server", "", "Rancher Server URL")
	flag.StringVar(&clusterName, "cluster-name", "", "Cluster Name")

	// Parse the command-line flags
	flag.Parse()

	fmt.Println("Cluster Name:", clusterName)
	fmt.Println("Create Kubeconfig:", createKubeconfig)
	fmt.Println("Create Namespace:", createNamespace)
	fmt.Println("Create Project:", createProject)
	fmt.Println("Filter Clusters by Type:", clusterType)
	fmt.Println("Filter Clusters by Label:", clusterLabels)
	fmt.Println("Kubeconfig File:", kubeconfigFile)
	fmt.Println("Namespace:", namespace)
	fmt.Println("Project Name:", projectName)
	fmt.Println("Rancher Access Key:", rancherAccessKey)
	fmt.Println("Rancher Secret Key:", rancherSecretKey)
	fmt.Println("Rancher Server URL:", rancherServer)
}

func printHelp() {
	fmt.Println("Usage: ./rancher-projects [options]")
	fmt.Println("Options:")
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Printf("  --%s %s\n", f.Name, f.Usage)
	})
}

func loadEnvVars() {
	rancherServer = os.Getenv("RANCHER_SERVER")
	rancherAccessKey = os.Getenv("RANCHER_ACCESS_KEY")
	rancherSecretKey = os.Getenv("RANCHER_SECRET_KEY")
	clusterName = os.Getenv("CLUSTER_NAME")
	projectName = os.Getenv("PROJECT_NAME")
	createProject = os.Getenv("CREATE_PROJECT") == "true"
	namespace = os.Getenv("NAMESPACE")
	createNamespace = os.Getenv("CREATE_NAMESPACE") == "true"
	clusterType = os.Getenv("CLUSTER_TYPE")
	clusterLabels = os.Getenv("CLUSTER_LABELS")
	kubeconfig = os.Getenv("KUBECONFIG")
	kubeconfigDir = os.Getenv("KUBECONFIG_DIR")
	kubeconfigPrefix = os.Getenv("KUBECONFIG_PREFIX")
	debug = os.Getenv("DEBUG") == "true"

	// Set default values if not provided
	if rancherServer == "" {
		rancherServer = "https://rancher.example.com"
	}
	if kubeconfigDir == "" {
		kubeconfigDir, _ = os.Getwd()
	}
	if kubeconfigPrefix == "" {
		kubeconfigPrefix = ""
	}
}

func verifySettings() {
	fmt.Println("Verifying settings...")

	if kubeconfig == "" {
		if debug {
			fmt.Println("Using Kubeconfig DIR")
		}

		if kubeconfigDir == "" {
			kubeconfigDir, _ = os.Getwd()
			if debug {
				fmt.Println("Defaulting to pwd")
			}
		} else {
			if debug {
				fmt.Println("Making Kubeconfig DIR")
			}
			err := os.MkdirAll(kubeconfigDir, os.ModePerm)
			if err != nil {
				fmt.Println("Kubeconfig directory does not exist. Please create it and try again.")
				os.Exit(1)
			}
		}
	}

	requiredVars := []string{
		"CLUSTER_NAME",
		"CLUSTER_TYPE",
		"CLUSTER_LABELS",
		"CATTLE_SERVER",
		"CATTLE_ACCESS_KEY",
		"CATTLE_SECRET_KEY",
	}

	for _, varName := range requiredVars {
		value := os.Getenv(varName)
		if value == "" {
			fmt.Printf("%s is required. Please specify it and try again.\n", varName)
			os.Exit(1)
		}
	}

	if debug {
		fmt.Println("Dumping options")
		fmt.Println("CLUSTER_NAME:", clusterName)
		fmt.Println("CLUSTER_TYPE:", clusterType)
		fmt.Println("CLUSTER_LABELS:", clusterLabels)
		fmt.Println("CATTLE_SERVER:", rancherServer)
		fmt.Println("CATTLE_ACCESS_KEY:", rancherAccessKey)
		fmt.Println("CATTLE_SECRET_KEY:", rancherSecretKey)
		fmt.Println("CREATE_KUBECONFIG:", createKubeconfig)
		fmt.Println("KUBECONFIG:", kubeconfig)
		fmt.Println("KUBECONFIG_DIR:", kubeconfigDir)
		fmt.Println("KUBECONFIG_PREFIX:", kubeconfigPrefix)
	}
}

func verifyAccess() {
	fmt.Println("Verifying access to Rancher server...")

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	url := fmt.Sprintf("%s/v3/", rancherServer)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to authenticate to %s\n", rancherServer)
		os.Exit(2)
	}

	fmt.Printf("Successfully authenticated to %s\n", rancherServer)
}

func verifyCluster() {
	fmt.Printf("Verifying cluster %s...\n", clusterName)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", rancherServer, clusterName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to find cluster %s\n", clusterName)
		os.Exit(2)
	}

	fmt.Printf("Successfully found cluster %s\n", clusterName)
}

func getClusterID() string {
	fmt.Printf("Getting cluster ID for %s...\n", clusterName)

	url := fmt.Sprintf("%s/v3/clusters?name=%s", rancherServer, clusterName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to get cluster ID")
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	if len(data.Data) == 0 {
		log.Fatal("Failed to find cluster ID")
	}

	clusterID := data.Data[0].ID
	fmt.Printf("Cluster ID: %s\n", clusterID)
	fmt.Println("Successfully got cluster ID")
	return clusterID
}

func verifyProject(clusterID, projectName string) {
	fmt.Println("Verifying project", projectName, "...")
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", rancherServer, clusterID, projectName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to verify project %s: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to verify project %s. Status code: %d", projectName, resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var projects []map[string]interface{}
	err = json.Unmarshal(body, &projects)
	if err != nil {
		log.Fatalf("Failed to parse project data: %v", err)
	}

	if len(projects) == 0 {
		log.Fatalf("Failed to find project %s", projectName)
	}

	fmt.Println("Successfully found project", projectName)
}

func performCreateProject() {
	fmt.Println("Checking if project", projectName, "exists...")
	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", rancherServer, clusterID, projectName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to check if project %s exists: %v", projectName, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Project", projectName, "already exists")
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Println("Creating project", projectName, "...")
		createProjectURL := fmt.Sprintf("%s/v3/projects", rancherServer)
		reqBody := fmt.Sprintf(`{"type": "project", "name": "%s", "clusterId": "%s"}`, projectName, clusterID)

		req, _ = http.NewRequest("POST", createProjectURL, strings.NewReader(reqBody))
		req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("Failed to create project %s: %v", projectName, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			log.Fatalf("Failed to create project %s. Status code: %d", projectName, resp.StatusCode)
		}

		fmt.Println("Successfully created project", projectName)
	} else {
		log.Fatalf("Failed to check project %s. Status code: %d", projectName, resp.StatusCode)
	}
}

func verifyNamespace(clusterID, namespace string) {
	fmt.Println("Verifying namespace", namespace, "...")
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", rancherServer, clusterID, namespace)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to verify namespace %s: %v", namespace, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Fatalf("Failed to find namespace %s", namespace)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to verify namespace %s. Status code: %d", namespace, resp.StatusCode)
	}

	fmt.Println("Successfully found namespace", namespace)
}

func performCreateNamespace() {
	fmt.Println("Checking if namespace", namespace, "exists...")
	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces", rancherServer, clusterID)
	reqBody := fmt.Sprintf(`{"type": "namespace", "metadata": {"name": "%s"}}`, namespace)

	req, _ := http.NewRequest("POST", url, strings.NewReader(reqBody))
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to check if namespace %s exists: %v", namespace, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Successfully created namespace", namespace)
		fmt.Println("Sleeping for 5 seconds to allow namespace to settle...")
		time.Sleep(5 * time.Second)
	} else if resp.StatusCode == http.StatusConflict {
		fmt.Println("Namespace", namespace, "already exists")
	} else {
		log.Fatalf("Failed to create namespace %s. Status code: %d", namespace, resp.StatusCode)
	}
}

func getClusterType(clusterName string) string {
	fmt.Println("Getting cluster type...")
	url := fmt.Sprintf("%s/v3/clusters?name=%s", rancherServer, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster type: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster type. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster type: %v", err)
	}

	if len(clusters) > 0 {
		return clusters[0]["provider"].(string)
	}

	return ""
}

func filterByClusterLabel(clusterName, keyPair string) bool {
	fmt.Println("Filtering by cluster label")
	fmt.Println("Cluster name:", clusterName)
	fmt.Println("Label:", keyPair)

	parts := strings.Split(keyPair, "=")
	if len(parts) != 2 {
		log.Fatalf("Invalid key-value pair: %s", keyPair)
	}

	key := parts[0]
	value := parts[1]

	url := fmt.Sprintf("%s/v3/clusters?name=%s", rancherServer, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
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
		labelValue, exists := labels[key]
		if exists && labelValue == value {
			return true
		}
	}

	return false
}

func getClusterStatus(clusterName string) string {
	fmt.Println("Getting cluster status...")
	url := fmt.Sprintf("%s/v3/clusters?name=%s", rancherServer, clusterName)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster status: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster status. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster status: %v", err)
	}

	if len(clusters) > 0 {
		return clusters[0]["state"].(string)
	}

	return ""
}

func getProjectInfo() string {
	fmt.Println("Getting project info...")

	url := fmt.Sprintf("%s/v3/projects?clusterId=%s&name=%s", rancherServer, clusterID, projectName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to get project info")
	}

	var data struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	if len(data.Data) == 0 {
		log.Fatal("Failed to find project info")
	}

	projectID := data.Data[0].ID
	fmt.Printf("Project ID: %s\n", projectID)
	fmt.Println("Successfully got project info")
	return projectID
}

func assignNamespaceToProject(projectID, namespace string) {
	fmt.Printf("Assigning namespace %s to project %s...\n", namespace, projectName)

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", rancherServer, clusterID, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to assign namespace to project")
	}

	var namespaceData struct {
		Metadata struct {
			Annotations struct {
				ProjectID string `json:"field.cattle.io/projectId"`
			} `json:"annotations"`
		} `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&namespaceData)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	namespaceData.Metadata.Annotations.ProjectID = projectID

	payload, err := json.Marshal(namespaceData)
	if err != nil {
		log.Fatal("Failed to marshal JSON payload:", err)
	}

	url = fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", rancherServer, clusterID, namespace)
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to assign namespace to project")
	}

	fmt.Printf("Successfully assigned namespace %s to project %s\n", namespace, projectName)
}

func verifyProjectAssignment(projectID, namespace string) {
	fmt.Println("Verifying project assignment...")

	url := fmt.Sprintf("%s/k8s/clusters/%s/v1/namespaces/%s", rancherServer, clusterID, namespace)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to send HTTP request:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Failed to verify project assignment")
	}

	var namespaceData struct {
		Metadata struct {
			Annotations struct {
				ProjectID string `json:"field.cattle.io/projectId"`
			} `json:"annotations"`
		} `json:"metadata"`
	}

	err = json.NewDecoder(resp.Body).Decode(&namespaceData)
	if err != nil {
		log.Fatal("Failed to decode JSON response:", err)
	}

	if namespaceData.Metadata.Annotations.ProjectID != projectID {
		log.Fatal("Failed to verify project assignment")
	}

	fmt.Println("Successfully verified project assignment")
}

func generateKubeconfig(kubeconfigFile, clusterID string) {
	fmt.Println("Generating kubeconfig...")

	url := fmt.Sprintf("%s/v3/clusters/%s?action=generateKubeconfig", rancherServer, clusterID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Fatal("Failed to create HTTP request:", err)
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", rancherAccessKey, rancherSecretKey)))
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

func getAllClusterIDs() {
	fmt.Println("Getting all cluster IDs...")
	url := fmt.Sprintf("%s/v3/clusters", rancherServer)
	req, _ := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(rancherAccessKey, rancherSecretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to retrieve cluster IDs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed to retrieve cluster IDs. Status code: %d", resp.StatusCode)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var clusters []map[string]interface{}
	err = json.Unmarshal(body, &clusters)
	if err != nil {
		log.Fatalf("Failed to parse cluster IDs: %v", err)
	}

	clusterIDs = make([]string, len(clusters))
	for i, cluster := range clusters {
		clusterName := cluster["name"].(string)
		clusterID := cluster["id"].(string)
		clusterIDs[i] = fmt.Sprintf("%s:%s", clusterName, clusterID)
	}

	fmt.Println("Cluster IDs:", clusterIDs)
}

func main() {
	if showHelp {
		printHelp()
		return
	}

	loadEnvVars()
	fmt.Println("Cluster Name:", clusterName)
	fmt.Println("Create Kubeconfig:", createKubeconfig)
	fmt.Println("Create Namespace:", createNamespace)
	fmt.Println("Create Project:", createProject)
	fmt.Println("Filter Clusters by Type:", getClustersByType)
	fmt.Println("Filter Clusters by Label:", getClustersByLabel)
	fmt.Println("Kubeconfig File:", kubeconfig)
	fmt.Println("Namespace:", namespace)
	fmt.Println("Project Name:", projectName)
	fmt.Println("Rancher Access Key:", rancherAccessKey)
	fmt.Println("Rancher Secret Key:", rancherSecretKey)
	fmt.Println("Rancher Server URL:", rancherServer)
	verifySettings()

	verifyAccess()

	if clusterType == "" && clusterLabels == "" {
		verifyCluster()
		getClusterID()

		if projectName != "" {
			if createProject {
				performCreateProject()
			}

			verifyProject(clusterID, projectName)
			getProjectInfo()

			if createNamespace {
				performCreateNamespace()
			} else {
				verifyNamespace(clusterID, namespace)
			}

			assignNamespaceToProject(clusterID, namespace)
			verifyProjectAssignment(clusterID, namespace)
		}

		if createKubeconfig {
			generateKubeconfig(kubeconfigFile, clusterID)
		}
	} else {
		getAllClusterIDs()

		keyPairs := strings.Split(clusterLabels, ",")
		keyPairCount := len(keyPairs)

		for _, clusterID := range clusterIDs {
			clusterName := strings.Split(clusterID, ":")[0]
			clusterID := strings.Split(clusterID, ":")[1]

			fmt.Println("Checking if cluster is Active...")
			getClusterStatus(clusterName)

			if clusterStatus == "active" {
				fmt.Println("Cluster is Active")

				if clusterType != "" {
					fmt.Println("Checking cluster type...")
					clusterProvider := getClusterType(clusterName)

					if clusterType == clusterProvider {
						fmt.Println("Cluster type match found")
						generateKubeconfig(clusterName, clusterID)
					}
				} else if clusterLabels != "" {
					found := 0
					foundAll := 0
					labelCount := 0

					for _, keyPair := range keyPairs {
						fmt.Println("Checking label", keyPair)
						filterByClusterLabel(clusterName, keyPair)
					}

					if labelCount == keyPairCount {
						fmt.Println("Found all labels")
						foundAll = 1
					} else {
						fmt.Println("Label count mismatch")
						foundAll = 0
					}

					if found == 1 || foundAll == 1 {
						fmt.Println("Matching label or no label set")
						generateKubeconfig(clusterName, clusterID)
					} else {
						fmt.Println("Skipping cluster due to missing label")
					}
				}
			} else {
				fmt.Println("Skipping Cluster because it is not Active")
			}

			fmt.Println("##################################################################################")
		}
	}
}
