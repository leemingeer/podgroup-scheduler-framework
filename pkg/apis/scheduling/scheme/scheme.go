package scheme

import (
	schedv1 "github.com/leemingeer/podgroup-scheduler-framework/pkg/apis/scheduling/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func init() {
	AddToScheme(scheme.Scheme)
}

// AddToScheme builds the kubescheduler scheme using all known versions of the kubescheduler api.
func AddToScheme(scheme *runtime.Scheme) {
	utilruntime.Must(schedv1.AddToScheme(scheme))
}