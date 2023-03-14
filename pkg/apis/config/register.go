package config

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	schedscheme "k8s.io/kubernetes/pkg/scheduler/apis/config"
)

// GroupName is the group name used in this package
const GroupName = "kubescheduler.config.k8s.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}

var (
	// reuse kube-scheduler scheme builder obj as local
	localSchemeBuilder = &schedscheme.SchemeBuilder
	// AddToScheme is a global function that registers this API group & version to a scheme
	// 外部通过这个方法，完成注册
	AddToScheme = localSchemeBuilder.AddToScheme
)

// addKnownTypes registers known types to the given scheme
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CoschedulingArgs{},
	)
	return nil
}

func init() {
	// We only register manually written functions here. The registration of the
	// generated functions takes place in the generated files. The separation
	// makes the code compile even when the generated files are missing.
	localSchemeBuilder.Register(addKnownTypes)
}



//// SchemeGroupVersion is group version used to register these objects
//var (
//	localSchemeBuilder = &schedscheme.SchemeBuilder
//	GroupVersion = schema.GroupVersion{Group: GroupName, Version: runtime.APIVersionInternal}
//)
//
//func init(){
//	Scheme.AddKnownTypes(GroupVersion, &CoschedulingArgs{})
//}