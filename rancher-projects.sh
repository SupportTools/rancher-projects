#!/bin/bash

usage() {
    echo "Usage: $0"
    echo "  --cluster-name              Cluster name                                 Example: rke2-dev"
    echo "  --project-name              Project name                                 Example: MyProject"
    echo "  --create-project            Creates project if it doesn't exist"
    echo "  --namespace                 Namespace)                                   Example: my-namespace"
    echo "  --create-namespace          Creates namespace if it doesn't exist"
    echo "  --rancher-server            Rancher server                               Example: https://rancher.dev.local"
    echo "  --rancher-access-key        Rancher Access Key                           Example: token-abcdef" 
    echo "  --rancher-secret-key        Rancher Secret Key                           Example: abcdefghijklmnopqrstuvwxyz"
    echo "  --create-kubeconfig         Generates kubeconfig file for the cluster"
    echo "  --get-clusters-by-type      Returns a list of clusters by type           Example: rke2"
    echo "  --kubeconfig                Overrides the kubeconfig file name           Default: rancher-projects-kubeconfig"
    echo "  --kubeconfig-dir            Overrides the kubeconfig file directory      Default: Current directory"
    echo "  --kubeconfig-prefix         Overrides the kubeconfig file prefix         Default: <blank>"
    exit 1; }

CREATE_PROJECT="false"
CREATE_NAMESPACE="false"
CREATE_KUBECONFIG="false"
KUBECONFIG_FILE="rancher-projects-kubeconfig"

POSITIONAL_ARGS=()

while [[ $# -gt 0 ]]; do
case $1 in
    -c|--cluster-name)
    CLUSTER_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    -p|--project-name)
    PROJECT_NAME="$2"
    shift # past argument
    shift # past value
    ;;
    --create-project)
    CREATE_PROJECT="true"
    shift # past argument
    shift # past value
    ;;
    -n|--namespace)
    NAMESPACE="$2"
    shift # past argument
    shift # past value
    ;;
    --create-namespace)
    CREATE_NAMESPACE="true"
    shift # past argument
    shift # past value
    ;;
    -s|--rancher-server)
    CATTLE_SERVER="$2"
    shift # past argument
    shift # past value
    ;;
    -A|--rancher-access-key)
    CATTLE_ACCESS_KEY="$2"
    shift # past argument
    shift # past value
    ;;
    -S|--rancher-secret-key)
    CATTLE_SECRET_KEY="$2"
    shift # past argument
    shift # past value
    ;;
    --create-kubeconfig)
    CREATE_KUBECONFIG="true"
    shift # past argument
    shift # past value
    ;;
    --get-clusters-by-type)
    CLUSTER_TYPE="$2"
    shift # past argument
    shift # past value
    ;;
    --kubeconfig)
    KUBECONFIG="$2"
    shift # past argument
    shift # past value
    ;;
    --kubeconfig-dir)
    KUBECONFIG_DIR="$2"
    shift # past argument
    shift # past value
    ;;
    --kubeconfig-prexfix)
    KUBECONFIG_PREFIX="$2"
    shift # past argument
    shift # past value
    ;;
    -d|--debug)
    DEBUG="true"
    shift # past argument
    shift # past value
    ;;
    -h|--help)
    usage
    shift # past argument
    shift # past value
    ;;
    -*|--*)
    echo "Unknown option $1"
    exit 1
    ;;
    *)
    POSITIONAL_ARGS+=("$1") # save positional arg
    shift # past argument
    ;;
esac
done
set -- "${POSITIONAL_ARGS[@]}" # restore positional parameters

verify-tools() {
    echo "Verifying tools..."
    if ! command -v jq >/dev/null 2>&1; then
        echo "jq is not installed. Please install jq and try again."
        exit 1
    fi
    if ! command -v curl >/dev/null 2>&1; then
        echo "curl is not installed. Please install curl and try again."
        exit 1
    fi
}

verify-settings() {
    echo "Verifying settings..."
    if [ -z "${KUBECONFIG}" ]; then
        if [ -z "${KUBECONFIG_DIR}" ]; then
            KUBECONFIG_DIR="$(pwd)"
        else
            mkdir -p "${KUBECONFIG_DIR}"
            if [ ! -d "${KUBECONFIG_DIR}" ]; then
                echo "Kubeconfig directory does not exist. Please create it and try again."
                exit 1
            fi
        fi
    fi
    if [[ -z ${CLUSTER_TYPE} ]]; then
        if [[ -z $CLUSTER_NAME ]] || [[ -z $PROJECT_NAME ]] || [[ -z $NAMESPACE ]] || [[ -z $CATTLE_SERVER ]] || [[ -z $CATTLE_ACCESS_KEY ]] || [[ -z $CATTLE_SECRET_KEY ]]; then
            usage
            exit 1
        fi
    else
        if [[ -z $CATTLE_SERVER ]] || [[ -z $CATTLE_ACCESS_KEY ]] || [[ -z $CATTLE_SECRET_KEY ]]; then
            usage
            exit 1
        fi
    fi
    if [ "${DEBUG}" == "true" ]; then
        set -x
        echo "Dumping options"
        echo "CLUSTER_NAME: ${CLUSTER_NAME}"
        echo "PROJECT_NAME: ${PROJECT_NAME}"
        echo "CREATE_PROJECT: ${CREATE_PROJECT}"
        echo "NAMESPACE: ${NAMESPACE}"
        echo "CREATE_NAMESPACE: ${CREATE_NAMESPACE}"
        echo "CATTLE_SERVER: ${CATTLE_SERVER}"
        echo "CATTLE_ACCESS_KEY: ${CATTLE_ACCESS_KEY}"
        echo "CATTLE_SECRET_KEY: ${CATTLE_SECRET_KEY}"
        echo "CREATE_KUBECONFIG: ${CREATE_KUBECONFIG}"
        echo "CLUSTER_TYPE: ${CLUSTER_TYPE}"
        echo "KUBECONFIG: ${KUBECONFIG}"
        echo "KUBECONFIG_DIR: ${KUBECONFIG_DIR}"
        echo "KUBECONFIG_PREFIX: ${KUBECONFIG_PREFIX}"
        echo "DEBUG: ${DEBUG}"
    fi

}

verify-access() {
    echo "Verifying access to Rancher server..."
    output=`curl -H 'content-type: application/json' -k -s -o /dev/null -w "%{http_code}" "${CATTLE_SERVER}/v3/" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}"`
    if [ $output -ne 200 ]; then
        echo "Failed to authenticate to ${CATTLE_SERVER}"
        exit 2
    fi
    echo "Successfully authenticated to ${CATTLE_SERVER}"
}

verify-cluster() {
    echo "Verifying cluster ${CLUSTER_NAME}..."
    output=`curl -H 'content-type: application/json' -k -s -o /dev/null -w "%{http_code}"  -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" "${CATTLE_SERVER}/v3/clusters?name=${CLUSTER_NAME}"
    if [ $output -ne 200 ]; then
        echo "Failed to find cluster ${CLUSTER_NAME}"
        exit 2
    fi
    echo "Successfully found cluster ${CLUSTER_NAME}"
}

verify-project() {
    echo "Verifying project ${PROJECT_NAME}..."
    PROJECT_DATA=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/projects?name=${PROJECT_NAME}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.data[0]'`
    if [[ "${PROJECT_DATA}" == "null" ]]; then
        echo "Failed to find project ${PROJECT_NAME}"
        exit 2
    fi
    echo "Successfully found project ${PROJECT_NAME}"
}

create-project() {
    echo "Checking if project ${PROJECT_NAME} exists..."
    PROJECT_DATA=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/projects?clusterId=${CLUSTER_ID}&name=${PROJECT_NAME}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.data[0]'`
    if [[ "${PROJECT_DATA}" == "null" ]]; then
        echo "Creating project ${PROJECT_NAME}..."
        curl -X POST -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/projects" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" -d "{\"type\":\"project\", \"name\":\"${PROJECT_NAME}\", \"clusterId\":\"${CLUSTER_ID}\"}" > /dev/null
        if [ $? -ne 0 ]; then
            echo "Failed to create project ${PROJECT_NAME}"
            exit 2
        fi
        echo "Successfully created project ${PROJECT_NAME}"
    else
        echo "Project ${PROJECT_NAME} already exists"
    fi
}

verify-namespace() {
    echo "Verifying namespace ${NAMESPACE}..."
    NAMESPACE_DATA=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces/${NAMESPACE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq .code | tr -d '"'`
    if [[ "${NAMESPACE_DATA}" == "NotFound" ]]; then
        echo "Failed to find namespace ${NAMESPACE}"
        exit 2
    fi
    echo "Successfully found namespace ${NAMESPACE}"
}

create-namespace() {
    echo "Checking if namespace ${NAMESPACE} exists..."
    NAMESPACE_DATA=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces/${NAMESPACE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq .code | tr -d '"'`
    if [[ "${NAMESPACE_DATA}" == "NotFound" ]]; then
        echo "Creating namespace ${NAMESPACE}..."
        curl -X POST -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" -d "{\"type\":\"namespace\", \"metadata\": {\"name\":\"${NAMESPACE}\"}}" > /dev/null
        verify-namespace
        echo "Successfully created namespace ${NAMESPACE}"
        echo "Sleep for 5 seconds to allow namespace to settle..."
        sleep 5
    else
        echo "Namespace ${NAMESPACE} already exists"
    fi
}

get-all-cluster-ids() {
    echo "Getting all cluster ids..."
    echo "CLUSTER_TYPE: ${CLUSTER_TYPE}"
    CLUSTER_IDS=$(curl -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/clusters?driver=${CLUSTER_TYPE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.data[] | .name + ":" + .id')
    if [ $? -ne 0 ]; then
        echo "Failed to get cluster id"
        exit 2
    fi
    echo "Successfully got cluster ids"
}

get-cluster-id() {
    echo "Getting cluster id..."
    CLUSTER_ID=$(curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/clusters?name=${CLUSTER_NAME}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.data[0].id')
    if [ $? -ne 0 ]; then
        echo "Failed to get cluster id"
        exit 2
    fi
    echo "Cluster id: ${CLUSTER_ID}"
    echo "Successfully got cluster id"
}

get-project-info() {
    echo "Getting project info..."
    PROJECT_ID=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/projects?clusterId=${CLUSTER_ID}&name=${PROJECT_NAME}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.data[0].id'`
    if [ $? -ne 0 ]; then
        echo "Failed to get project info"
        exit 2
    fi
    echo "Project id: ${PROJECT_ID}"
}

assign-namespace-to-project() {
    echo "Assigning namespace ${NAMESPACE} to project ${PROJECT_NAME}..."
    echo "Collecting namespace details..."
    PROJECT=`echo ${PROJECT_ID} | awk -F ':' '{print $2}'`
    echo "Project long: ${PROJECT_ID}"
    echo "Project short: ${PROJECT}"
    NAMESPACE_DATA=`curl  -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces/${NAMESPACE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq`
    NAMESPACE_POST=`echo ${NAMESPACE_DATA} | jq --arg PROJECT_ID "${PROJECT_ID}" '.metadata.annotations."field.cattle.io/projectId" = $PROJECT_ID' | jq --arg PROJECT "${PROJECT}" '.metadata.labels."field.cattle.io/projectId" = $PROJECT'`
    echo "Updating namespace..."
    curl -X PUT -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces/${NAMESPACE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" -d "${NAMESPACE_POST}" > /dev/null
    if [ $? -ne 0 ]; then
        echo "Failed to assign namespace ${NAMESPACE} to project ${PROJECT_NAME}"
        exit 2
    fi
    echo "Successfully assigned namespace ${NAMESPACE} to project ${PROJECT_NAME}"
}

verify-project-assignment() {
    echo "Verifying project assignment..."
    NAMESPACE_DATA=`curl -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/k8s/clusters/${CLUSTER_ID}/v1/namespaces/${NAMESPACE}" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq .metadata.annotations | grep "field.cattle.io/projectId" | awk '{print $2}' | tr -d '", '`
    if [[ "${NAMESPACE_DATA}" != "${PROJECT_ID}" ]]; then
        echo "Failed to verify project assignment"
        exit 2
    fi
    echo "Successfully verified project assignment"
}

generate-kubeconfig() {
    echo "Generating kubeconfig..."
    if [[ -z "${KUBECONFIG}" ]]; then
        KUBECONFIG_FILE=$1
        CLUSTER_ID=$2
    else
        KUBECONFIG_FILE=${KUBECONFIG}
    fi
    if [[ ! -z ${KUBECONFIG_PREFIX} ]]
    then
        KUBECONFIG_FILE="${KUBECONFIG_PREFIX}-${KUBECONFIG_FILE}"
    fi
    echo "Kubeconfig file: ${KUBECONFIG_FILE}"
    echo "Cluster id: ${CLUSTER_ID}"
    curl -X POST -H 'content-type: application/json' -k -s "${CATTLE_SERVER}/v3/clusters/${CLUSTER_ID}?action=generateKubeconfig" -u "${CATTLE_ACCESS_KEY}:${CATTLE_SECRET_KEY}" | jq -r '.config' > ${KUBECONFIG_DIR}/${KUBECONFIG_FILE}
}

verify-tools
verify-settings
verify-access
if [ -z ${CLUSTER_TYPE} ]; then
    verify-cluster
    get-cluster-id
    if [ "${CREATE_PROJECT}" == "true" ]; then
        create-project
    fi
    verify-project
    get-project-info
    if [ "${CREATE_NAMESPACE}" == "true" ]; then
        create-namespace
    else
        verify-namespace
    fi
    assign-namespace-to-project
    verify-project-assignment
    if [ "${CREATE_KUBECONFIG}" == "true" ]; then
        generate-kubeconfig
    fi
else
    get-all-cluster-ids
    for CLUSTER_ID in ${CLUSTER_IDS}; do
        cluster_name=`echo ${CLUSTER_ID} | awk -F ':' '{print $1}'`
        cluster_id=`echo ${CLUSTER_ID} | awk -F ':' '{print $2}'`
        generate-kubeconfig ${cluster_name} ${cluster_id}
    done
fi