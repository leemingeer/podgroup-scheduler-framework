package util

import "fmt"

const (
	// PodGroupLabel is the default label of coscheduling
	PodGroupLabel = "pod-group.scheduling.ming.io"
)

var (
	// ErrorNotMatched means pod does not match coscheduling
	ErrorNotMatched = fmt.Errorf("not match coscheduling")
	// ErrorWaiting means pod number does not match the min pods required
	ErrorWaiting = fmt.Errorf("waiting")
	// ErrorResourceNotEnough means cluster resource is not enough, mainly used in Pre-Filter
	ErrorResourceNotEnough = fmt.Errorf("resource not enough")
)
