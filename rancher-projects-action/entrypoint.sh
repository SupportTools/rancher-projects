#!/bin/bash

set -e

RANCHER_SERVER=$1
RANCHER_ACCESS_KEY=$2
RANCHER_SECRET_KEY=$3
CLUSTER_NAME=$4
PROJECT_NAME=$5
NAMESPACE=$6
CREATE_PROJECT=$7
CREATE_NAMESPACE=$8
CREATE_KUBECONFIG=$9
KUBECONFIG=${10}

# Download the rancher-projects script
curl -o /usr/local/bin/rancher-projects https://raw.githubusercontent.com/SupportTools/rancher-projects/main/rancher-projects.sh
chmod +x /usr/local/bin/rancher-projects

# Run the rancher-projects script with the provided inputs
bash /usr/local/bin/rancher-projects \
  --rancher-server "$RANCHER_SERVER" \
  --rancher-access-key "$RANCHER_ACCESS_KEY" \
  --rancher-secret-key "$RANCHER_SECRET_KEY" \
  --cluster-name "$CLUSTER_NAME" \
  --project-name "$PROJECT_NAME" \
  --namespace "$NAMESPACE" \
  --create-project "$CREATE_PROJECT" \
  --create-namespace "$CREATE_NAMESPACE" \
  --create-kubeconfig "$CREATE_KUBECONFIG" \
  --kubeconfig "$KUBECONFIG"
