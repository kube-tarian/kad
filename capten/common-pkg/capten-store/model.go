package captenstore

import (
	"strings"
	"time"

	"database/sql/driver"

	"github.com/google/uuid"
)

// StringArray represents a string array to handle TEXT[] data type
type StringArray []string

// GormDataType returns the data type of the field
func (StringArray) GormDataType() string {
	return "text[]"
}

// Value gets the value to store in the database
func (a StringArray) Value() (driver.Value, error) {
	var arr = "{" + a[0]
	for _, v := range a[1:] {
		arr += "," + v
	}
	arr += "}"
	return arr, nil
}

func (a StringArray) String() string {
	var arr = "{\"" + a[0]
	for _, v := range a[1:] {
		arr += "\",\"" + v
	}
	arr += "\"}"
	return arr
}

// Scan reads the value from the database
func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	input := strings.Trim(value.(string), "{}")
	parts := strings.Split(input, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	*a = parts
	return nil
}

type GitProject struct {
	ID             uuid.UUID   `json:"id" gorm:"column:id;primaryKey"`
	ProjectURL     string      `json:"project_url" gorm:"column:project_url"`
	Labels         StringArray `json:"labels" gorm:"column:labels;type:text[]"`
	LastUpdateTime time.Time   `json:"last_update_time column:last_update_time"`
}

func (GitProject) TableName() string {
	return "git_project"
}

type CloudProvider struct {
	ID             uuid.UUID   `json:"id" gorm:"column:id;type:uuid"`
	CloudType      string      `json:"cloud_type" gorm:"column:cloud_type"`
	Labels         StringArray `json:"labels" gorm:"column:labels;type:text[]"`
	LastUpdateTime time.Time   `json:"last_update_time" gorm:"column:last_update_time"`
}

func (CloudProvider) TableName() string {
	return "cloud_provider"
}

type ContainerRegistry struct {
	ID             uuid.UUID   `json:"id" gorm:"column:id;primaryKey"`
	RegistryURL    string      `json:"registry_url" gorm:"column:registry_url"`
	RegistryType   string      `json:"registry_type" gorm:"column:registry_type"`
	Labels         StringArray `json:"labels" gorm:"column:labels;type:text[]"`
	LastUpdateTime time.Time   `json:"last_update_time" gorm:"column:last_update_time"`
}

func (ContainerRegistry) TableName() string {
	return "container_registry"
}

type ClusterAppConfig struct {
	ReleaseName         string    `json:"release_name" gorm:"column:release_name;primaryKey"`
	AppName             string    `json:"app_name" gorm:"column:app_name"`
	PluginName          string    `json:"plugin_name" gorm:"column:plugin_name"`
	PluginStoreType     int       `json:"plugin_store_type" gorm:"column:plugin_store_type"`
	Category            string    `json:"category" gorm:"column:category"`
	Description         string    `json:"description" gorm:"column:description"`
	RepoURL             string    `json:"repo_url" gorm:"column:repo_url"`
	Version             string    `json:"version" gorm:"column:version"`
	Namespace           string    `json:"namespace" gorm:"column:namespace"`
	UIEndpoint          string    `json:"ui_endpoint" gorm:"column:ui_endpoint"`
	UIModuleEndpoint    string    `json:"ui_module_endpoint" gorm:"column:ui_module_endpoint"`
	APIEndpoint         string    `json:"api_endpoint" gorm:"column:api_endpoint"`
	DefaultApp          bool      `json:"default_app" gorm:"column:default_app"`
	PrivilegedNamespace bool      `json:"privileged_namespace" gorm:"column:privileged_namespace"`
	InstallStatus       string    `json:"install_status" gorm:"column:install_status"`
	Icon                []byte    `json:"icon" gorm:"column:icon"`
	OverrideValues      string    `json:"override_values" gorm:"column:override_values"`
	LaunchUIValues      string    `json:"launch_ui_values" gorm:"column:launch_ui_values"`
	TemplateValues      string    `json:"template_values" gorm:"column:template_values"`
	LastUpdateTime      time.Time `json:"last_update_time" gorm:"column:last_update_time"`
}

func (ClusterAppConfig) TableName() string {
	return "cluster_app_config"
}

type ClusterPluginConfig struct {
	PluginName          string      `json:"plugin_name" gorm:"column:plugin_name;primaryKey"`
	PluginStoreType     int         `json:"plugin_store_type" gorm:"column:plugin_store_type"`
	Capabilities        StringArray `json:"capabilities" gorm:"column:capabilities;type:text[]"`
	Category            string      `json:"category" gorm:"column:category"`
	Description         string      `json:"description" gorm:"column:description"`
	ChartName           string      `json:"chart_name" gorm:"column:chart_name"`
	ChartRepo           string      `json:"chart_repo" gorm:"column:chart_repo"`
	Version             string      `json:"version" gorm:"column:version"`
	Namespace           string      `json:"namespace" gorm:"column:namespace"`
	UIEndpoint          string      `json:"ui_endpoint" gorm:"column:ui_endpoint"`
	APIEndpoint         string      `json:"api_endpoint" gorm:"column:api_endpoint"`
	UIModuleEndpoint    string      `json:"ui_module_endpoint" gorm:"column:ui_module_endpoint"`
	DefaultApp          bool        `json:"default_app" gorm:"column:default_app"`
	PrivilegedNamespace bool        `json:"privileged_namespace" gorm:"column:privileged_namespace"`
	InstallStatus       string      `json:"install_status" gorm:"column:install_status"`
	Icon                []byte      `json:"icon" gorm:"column:icon"`
	OverrideValues      string      `json:"override_values" gorm:"column:override_values"`
	Values              string      `json:"values" gorm:"column:values"`
	LastUpdateTime      time.Time   `json:"last_update_time" gorm:"column:last_update_time"`
}

func (ClusterPluginConfig) TableName() string {
	return "cluster_plugin_config"
}

type ManagedCluster struct {
	ID                  uuid.UUID `json:"id" gorm:"column:id;primaryKey"`
	ClusterName         string    `json:"cluster_name" gorm:"column:cluster_name"`
	ClusterEndpoint     string    `json:"cluster_endpoint" gorm:"column:cluster_endpoint"`
	ClusterDeployStatus string    `json:"cluster_deploy_status" gorm:"column:cluster_deploy_status"`
	AppDeployStatus     string    `json:"app_deploy_status" gorm:"column:app_deploy_status"`
	LastUpdateTime      time.Time `json:"last_update_time" gorm:"column:last_update_time"`
}

func (ManagedCluster) TableName() string {
	return "managed_clusters"
}

type TektonProject struct {
	ID             int       `json:"id" gorm:"column:id;primaryKey"`
	GitProjectID   uuid.UUID `json:"git_project_id" gorm:"column:git_project_id"`
	GitProjectURL  string    `json:"git_project_url" gorm:"column:git_project_url"`
	Status         string    `json:"status" gorm:"column:status"`
	LastUpdateTime time.Time `json:"last_update_time" gorm:"column:last_update_time"`
}

func (TektonProject) TableName() string {
	return "tekton_project"
}

type CrossplaneProvider struct {
	ID              uuid.UUID `json:"id" gorm:"column:id;primaryKey"`
	CloudProviderID string    `json:"cloud_provider_id" gorm:"column:cloud_provider_id"`
	ProviderName    string    `json:"provider_name" gorm:"column:provider_name"`
	CloudType       string    `json:"cloud_type" gorm:"column:cloud_type"`
	Status          string    `json:"status" gorm:"column:status"`
	LastUpdateTime  time.Time `json:"last_update_time" gorm:"column:last_update_time"`
}

func (CrossplaneProvider) TableName() string {
	return "crossplane_provider"
}

type CrossplaneProject struct {
	ID             int       `json:"id"`
	GitProjectID   uuid.UUID `json:"git_project_id"`
	GitProjectURL  string    `json:"git_project_url"`
	Status         string    `json:"status"`
	LastUpdateTime time.Time `json:"last_update_time"`
}

func (CrossplaneProject) TableName() string {
	return "crossplane_project"
}
