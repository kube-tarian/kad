package model

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	providerNamePrefix       = "provider"
	CrossPlaneResource       = "crossplane"
	CrossPlaneClusterUpdate  = "crossplane-cluster-update"
	CrossPlaneProjectSync    = "crossplane-project-sync"
	CrossPlaneProjectDelete  = "crossplane-project-delete"
	CrossPlaneProviderUpdate = "crossplane-provider-update"
)

type CrossplaneProviderStatus string

const (
	CrossPlaneProviderOutofSynch    CrossplaneProviderStatus = "OutOfSynch"
	CrossPlaneProviderInSynch       CrossplaneProviderStatus = "InSynch"
	CrossPlaneProviderFailedToSynch CrossplaneProviderStatus = "FailedToSynch"
	CrossPlaneProviderReady         CrossplaneProviderStatus = "Ready"
	CrossPlaneProviderNotReady      CrossplaneProviderStatus = "NotReady"
)

type CrossplaneProvider struct {
	Id              string `json:"id,omitempty"`
	CloudType       string `json:"cloud_type,omitempty"`
	ProviderName    string `json:"provider_name,omitempty"`
	CloudProviderId string `json:"cloud_provider_id,omitempty"`
	Status          string `json:"status,omitempty"`
}

type CrossplaneProjectStatus string

const (
	CrossplaneProjectAvailable            CrossplaneProjectStatus = "available"
	CrossplaneProjectConfigured           CrossplaneProjectStatus = "configured"
	CrossplaneProjectConfigurationOngoing CrossplaneProjectStatus = "configuration-ongoing"
	CrossplaneProjectConfigurationFailed  CrossplaneProjectStatus = "configuration-failed"
)

type CrossplaneProject struct {
	Id             string `json:"id,omitempty"`
	GitProjectId   string `json:"git_project_id,omitempty"`
	GitProjectUrl  string `json:"git_project_url,omitempty"`
	Status         string `json:"status,omitempty"`
	LastUpdateTime string `json:"last_update_time,omitempty"`
}

func PrepareCrossplaneProviderName(providerType string) string {
	return fmt.Sprintf("%s-%s", providerNamePrefix, strings.ToLower(providerType))
}

// A ConditionType represents a condition a resource could be in.
type ConditionType string

// Condition types.
const (
	TypeHealthy   ConditionType = "Healthy"
	TypeReady     ConditionType = "Ready"
	TypeInstalled ConditionType = "Installed"
	TypeSynced    ConditionType = "Synced"
)

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// Reasons a resource is or is not ready.
const (
	ReasonAvailable   ConditionReason = "Available"
	ReasonUnavailable ConditionReason = "Unavailable"
	ReasonCreating    ConditionReason = "Creating"
	ReasonDeleting    ConditionReason = "Deleting"
)

// Reasons a resource is or is not synced.
const (
	ReasonReconcileSuccess ConditionReason = "ReconcileSuccess"
	ReasonReconcileError   ConditionReason = "ReconcileError"
	ReasonReconcilePaused  ConditionReason = "ReconcilePaused"
)

type ClusterClaimSpec struct {
	Id string `json:"id,omitempty"`
}

type ClusterClaimCondition struct {
	LastTransitionTime string `json:"lastTransitionTime,omitempty" protobuf:"bytes,1,opt,name=lastTransitionTime"`
	Reason             string `json:"reason,omitempty" protobuf:"bytes,2,opt,name=reason"`
	Status             string `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	Type               string `json:"type,omitempty" protobuf:"bytes,4,opt,name=type"`
}

type ClusterClaimStatus struct {
	Conditions         []ClusterClaimCondition `json:"conditions,omitempty" protobuf:"bytes,1,opt,name=conditions"`
	ControlPlaneStatus string                  `json:"controlPlaneStatus,omitempty"`
	NodePoolStatus     string                  `json:"nodePoolStatus,omitempty"`
}

type ClusterClaim struct {
	Metadata metav1.ObjectMeta  `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec     ClusterClaimSpec   `json:"spec,omitempty" protobuf:"bytes,1,opt,name=spec"`
	Status   ClusterClaimStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

type ClusterClaimList struct {
	Items []ClusterClaim `json:"items,omitempty" protobuf:"bytes,1,opt,name=items"`
}

type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status ProviderStatus `json:"status,omitempty"`
}

// A Condition that may apply to a resource.
type Condition struct {
	// Type of this condition. At most one of each condition type may apply to
	// a resource at any point in time.
	Type ConditionType `json:"type"`

	// Status of this condition; is it currently True, False, or Unknown?
	Status corev1.ConditionStatus `json:"status"`

	// LastTransitionTime is the last time this condition transitioned from one
	// status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// A Reason for this condition's last transition from one status to another.
	Reason ConditionReason `json:"reason"`

	// A Message containing details about this condition's last transition from
	// one status to another, if any.
	// +optional
	Message string `json:"message,omitempty"`
}

type ConditionedStatus struct {
	Conditions []Condition `json:"conditions,omitempty"`
}

// ProviderStatus represents the observed state of a Provider.
type ProviderStatus struct {
	ConditionedStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProviderList contains a list of Provider.
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Provider `json:"items"`
}
