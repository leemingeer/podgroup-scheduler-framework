package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeschedulerscheme "k8s.io/kubernetes/pkg/scheduler/apis/config/scheme"
)

// GroupName is the group name used in this package
const GroupName = "kubescheduler.config.k8s.io"


var (
	// Re-use the in-tree Scheme.
	Scheme = kubeschedulerscheme.Scheme
	GroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}
)

func init(){
	Scheme.AddKnownTypes(GroupVersion, &CoschedulingArgs{})
}