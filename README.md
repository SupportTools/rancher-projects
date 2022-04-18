# rancher-projects
The goal of this project is to provide a tool for creating projects and namespaces in Rancher then assigns a namespace to a project.

## Installation
```bash
sudo curl -o /usr/local/bin/rancher-projects https://raw.githubusercontent.com/SupportTools/rancher-projects/main/rancher-projects.sh
sudo chmod +x /usr/local/bin/rancher-projects
```

of

```bash
wget -O rancher-projects https://raw.githubusercontent.com/SupportTools/rancher-projects/main/rancher-projects.sh
chmod +x rancher-projects
sudo mv rancher-projects /usr/local/bin/
```

## Usage
```bash
bash run.sh \
--rancher-server "https://rancher.mattox.local" \
--rancher-access-key "token-abcde" \
--rancher-secret-key "123456789abcdefghijklmnopqrstuvwxyz" \
--cluster-name "MyCluster" \
--project-name "MyProject" \
--namespace "mynamespace"  \
--create-project true \
--create-namespace true \
--create-kubeconfig true \
--kubeconfig "rancher-projects-kubeconfig"
```

## Options
`--rancher-server` sets the Rancher Server. Note: This should include `https://` 
`--rancher-access-key` sets the Rancher Access Key. Note: This account should have permissions to list clusters, get/list/create/update projects and namespace.
`--rancher-secret-key` sets the Rancher Secret Key.
`--cluster-name` sets the cluster name in which the project will be created.
`--project-name` sets the project name to be created and/or assigned.
`--namespace` sets the namespace name to be created and/or assigned.
`--create-project` sets whether to create the project. (Optional) If project does not exist, it will be created.
`--create-namespace` sets whether to create the namespace. (Optional) If namespace does not exist, it will be created.
`--create-kubeconfig` sets whether to create a kubeconfig file. (Optional) If kubeconfig file does not exist, it will be created.
`--kubeconfig` sets the path to the kubeconfig file. (Optional) Default is rancher-projects-kubeconfig.
`--help` prints this help message.

## Examples
```bash
rancher-projects --rancher-server "https://rancher.mattox.local" --rancher-access-key "token-abcde" --rancher-secret-key "123456789abcdefghijklmnopqrstuvwxyz"  --cluster-name a0-rke2-devops --project-name "ClusterServices" --namespace "monitoring" --create-project true --create-namespace true
Verifying tools...
Verifying access to Rancher server...
Successfully authenticated to https://rancher.mattox.local
Verifying cluster a0-rke2-devops...
Successfully found cluster a0-rke2-devops
Getting cluster id...
Cluster id: c-m-9ldt7ts5
Successfully got cluster id
Checking if project ClusterServices exists...
Project ClusterServices already exists
Getting project info...
Project id: c-m-9ldt7ts5:p-bn9h5
Checking if namespace monitoring exists...
Namespace monitoring already exists
Assigning namespace monitoring to project ClusterServices...
Collecting namespace details...
Project long: c-m-9ldt7ts5:p-bn9h5
Project short: p-bn9h5
Updating namespace...
Successfully assigned namespace monitoring to project ClusterServices
Generating kubeconfig...
Kubeconfig: rancher-projects-kubeconfig
```