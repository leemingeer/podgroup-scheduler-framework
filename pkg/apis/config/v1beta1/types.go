package v1beta1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CoschedulingArgs defines the scheduling parameters for Coscheduling plugin.
type CoschedulingArgs struct {
	metav1.TypeMeta `json:",inline"`

	// PermitWaitingTime is the wait timeout in seconds.
	PermitWaitingTimeSeconds *int64 `json:"permitWaitingTimeSeconds,omitempty"`
	// DeniedPGExpirationTimeSeconds is the expiration time of the denied podgroup store.
	DeniedPGExpirationTimeSeconds *int64 `json:"deniedPGExpirationTimeSeconds,omitempty"`
	// KubeMaster is the url of api-server
	KubeMaster *string `json:"kubeMaster,omitempty"`
	// KubeConfigPath for scheduler
	KubeConfigPath *string `json:"kubeConfigPath,omitempty"`
}
