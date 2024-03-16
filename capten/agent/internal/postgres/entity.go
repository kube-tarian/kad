package postgres

import (
	"time"

	"github.com/google/uuid"
)

type GitProjects struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ProjectURL     string    `json:"project_url"`
	Labels         []string  `json:"labels" gorm:"type:text[]"`
	LastUpdateTime time.Time `json:"last_update_time"`
	UsedPlugins    []string  `json:"used_plugins" gorm:"type:text[]"`
}

type CloudProviders struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CloudType      string    `json:"cloud_type"`
	Labels         []string  `json:"labels" gorm:"type:text[]"`
	LastUpdateTime time.Time `json:"last_update_time"`
	UsedPlugins    []string  `json:"used_plugins" gorm:"type:text[]"`
}

type ContainerRegistry struct {
	ID             uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RegistryURL    string    `json:"registry_url"`
	RegistryType   string    `json:"registry_type"`
	Labels         []string  `json:"labels" gorm:"type:text[]"`
	LastUpdateTime time.Time `json:"last_update_time"`
	UsedPlugins    []string  `json:"used_plugins" gorm:"type:text[]"`
}
