package rancher

import (
	"time"

	"github.com/supporttools/rancher-projects/pkg/logging"
)

var (
	logger = logging.SetupLogging()
)

type RancherResponse struct {
	Type         string      `json:"type"`
	Links        Links       `json:"links"`
	CreateTypes  CreateTypes `json:"createTypes"`
	Actions      struct{}    `json:"actions"`
	Pagination   Pagination  `json:"pagination"`
	Sort         Sort        `json:"sort"`
	Filters      Filters     `json:"filters"`
	ResourceType string      `json:"resourceType"`
	Data         []Project   `json:"data"`
}

type Links struct {
	Self string `json:"self"`
}

type CreateTypes struct {
	Project string `json:"project"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type Sort struct {
	Order   string    `json:"order"`
	Reverse string    `json:"reverse"`
	Links   SortLinks `json:"links"`
}

type SortLinks struct {
	Description                 string `json:"description"`
	Name                        string `json:"name"`
	PodSecurityPolicyTemplateId string `json:"podSecurityPolicyTemplateId"`
	State                       string `json:"state"`
	Transitioning               string `json:"transitioning"`
	TransitioningMessage        string `json:"transitioningMessage"`
	Uuid                        string `json:"uuid"`
}

type Filters struct {
	ClusterId                   []ModifierValue `json:"clusterId"`
	Created                     interface{}     `json:"created"`
	CreatorId                   interface{}     `json:"creatorId"`
	Description                 interface{}     `json:"description"`
	EnableProjectMonitoring     interface{}     `json:"enableProjectMonitoring"`
	Id                          interface{}     `json:"id"`
	Name                        []ModifierValue `json:"name"`
	NamespaceId                 interface{}     `json:"namespaceId"`
	PodSecurityPolicyTemplateId interface{}     `json:"podSecurityPolicyTemplateId"`
	Removed                     interface{}     `json:"removed"`
	State                       interface{}     `json:"state"`
	Transitioning               interface{}     `json:"transitioning"`
	TransitioningMessage        interface{}     `json:"transitioningMessage"`
	Uuid                        interface{}     `json:"uuid"`
}

type ModifierValue struct {
	Modifier string `json:"modifier"`
	Value    string `json:"value"`
}

type Project struct {
	Actions                 ProjectActions    `json:"actions"`
	Annotations             map[string]string `json:"annotations"`
	BaseType                string            `json:"baseType"`
	ClusterId               string            `json:"clusterId"`
	Conditions              []Condition       `json:"conditions"`
	Created                 time.Time         `json:"created"`
	CreatedTS               int64             `json:"createdTS"`
	CreatorId               string            `json:"creatorId"`
	EnableProjectMonitoring bool              `json:"enableProjectMonitoring"`
	Id                      string            `json:"id"`
	Labels                  map[string]string `json:"labels"`
	Links                   ProjectLinks      `json:"links"`
	Name                    string            `json:"name"`
	NamespaceId             interface{}       `json:"namespaceId"`
	State                   string            `json:"state"`
	Transitioning           string            `json:"transitioning"`
	TransitioningMessage    string            `json:"transitioningMessage"`
	Type                    string            `json:"type"`
	Uuid                    string            `json:"uuid"`
}

type ProjectActions struct {
	EnableMonitoring             string `json:"enableMonitoring"`
	ExportYaml                   string `json:"exportYaml"`
	SetPodSecurityPolicyTemplate string `json:"setpodsecuritypolicytemplate"`
}

type Condition struct {
	LastUpdateTime string `json:"lastUpdateTime"`
	Status         string `json:"status"`
	Type           string `json:"type"`
}

type ProjectLinks struct {
	Alertmanagers                            string `json:"alertmanagers"`
	AppRevisions                             string `json:"appRevisions"`
	Apps                                     string `json:"apps"`
	BasicAuths                               string `json:"basicAuths"`
	Certificates                             string `json:"certificates"`
	ConfigMaps                               string `json:"configMaps"`
	CronJobs                                 string `json:"cronJobs"`
	DaemonSets                               string `json:"daemonSets"`
	Deployments                              string `json:"deployments"`
	DnsRecords                               string `json:"dnsRecords"`
	DockerCredentials                        string `json:"dockerCredentials"`
	HorizontalPodAutoscalers                 string `json:"horizontalPodAutoscalers"`
	Ingresses                                string `json:"ingresses"`
	Jobs                                     string `json:"jobs"`
	NamespacedBasicAuths                     string `json:"namespacedBasicAuths"`
	NamespacedCertificates                   string `json:"namespacedCertificates"`
	NamespacedDockerCredentials              string `json:"namespacedDockerCredentials"`
	NamespacedSecrets                        string `json:"namespacedSecrets"`
	NamespacedServiceAccountTokens           string `json:"namespacedServiceAccountTokens"`
	NamespacedSshAuths                       string `json:"namespacedSshAuths"`
	PersistentVolumeClaims                   string `json:"persistentVolumeClaims"`
	PodSecurityPolicyTemplateProjectBindings string `json:"podSecurityPolicyTemplateProjectBindings"`
	Pods                                     string `json:"pods"`
	ProjectAlertGroups                       string `json:"projectAlertGroups"`
	ProjectAlertRules                        string `json:"projectAlertRules"`
	ProjectAlerts                            string `json:"projectAlerts"`
	ProjectCatalogs                          string `json:"projectCatalogs"`
	ProjectMonitorGraphs                     string `json:"projectMonitorGraphs"`
	ProjectNetworkPolicies                   string `json:"projectNetworkPolicies"`
	ProjectRoleTemplateBindings              string `json:"projectRoleTemplateBindings"`
	PrometheusRules                          string `json:"prometheusRules"`
	Prometheuses                             string `json:"prometheuses"`
	Remove                                   string `json:"remove"`
	ReplicaSets                              string `json:"replicaSets"`
	ReplicationControllers                   string `json:"replicationControllers"`
	Secrets                                  string `json:"secrets"`
	Self                                     string `json:"self"`
	ServiceAccountTokens                     string `json:"serviceAccountTokens"`
	ServiceMonitors                          string `json:"serviceMonitors"`
	Services                                 string `json:"services"`
	SshAuths                                 string `json:"sshAuths"`
	StatefulSets                             string `json:"statefulSets"`
	Subscribe                                string `json:"subscribe"`
	Templates                                string `json:"templates"`
	Update                                   string `json:"update"`
	Workloads                                string `json:"workloads"`
}
