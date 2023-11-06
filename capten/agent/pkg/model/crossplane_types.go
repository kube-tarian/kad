package model

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	Conditions []ClusterClaimCondition `json:"conditions,omitempty" protobuf:"bytes,1,opt,name=conditions"`
}

type ClusterClaim struct {
	Metadata metav1.ObjectMeta  `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Spec     ClusterClaimSpec   `json:"spec,omitempty" protobuf:"bytes,1,opt,name=spec"`
	Status   ClusterClaimStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

type ClusterClaimList struct {
	Items []ClusterClaim `json:"items,omitempty" protobuf:"bytes,1,opt,name=items"`
}
