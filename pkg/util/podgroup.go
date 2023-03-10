package util

import (
	"encoding/json"
	"fmt"
	schedulingv1 "github.com/leemingeer/podgroup-scheduler-framework/pkg/apis/scheduling/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"time"
)
// GetPodGroupLabel get pod group from pod annotations
func GetPodGroupLabel(pod *v1.Pod) string {
	return pod.Labels[PodGroupLabel]
}

// DefaultWaitTime is 60s if ScheduleTimeoutSeconds is not specified.
const DefaultWaitTime = 60 * time.Second

// CreateMergePatch return patch generated from original and new interfaces
func CreateMergePatch(original, new interface{}) ([]byte, error) {
	pvByte, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	cloneByte, err := json.Marshal(new)
	if err != nil {
		return nil, err
	}
	patch, err := strategicpatch.CreateTwoWayMergePatch(pvByte, cloneByte, original)
	if err != nil {
		return nil, err
	}
	return patch, nil
}

// GetPodGroupFullName get namespaced group name from pod annotations
func GetPodGroupFullName(pod *v1.Pod) string {
	pgName := GetPodGroupLabel(pod)
	if len(pgName) == 0 {
		return ""
	}
	return fmt.Sprintf("%v/%v", pod.Namespace, pgName)
}

// GetWaitTimeDuration returns a wait timeout based on the following precedences:
// 1. spec.scheduleTimeoutSeconds of the given pg, if specified
// 2. given scheduleTimeout, if not nil
// 3. fall back to DefaultWaitTime
func GetWaitTimeDuration(pg *schedulingv1.PodGroup, scheduleTimeout *time.Duration) time.Duration {
	if pg != nil && pg.Spec.ScheduleTimeoutSeconds != nil {
		return time.Duration(*pg.Spec.ScheduleTimeoutSeconds) * time.Second
	}
	if scheduleTimeout != nil && *scheduleTimeout != 0 {
		return *scheduleTimeout
	}
	return DefaultWaitTime
}
