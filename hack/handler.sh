#!/bin/sh

# generate deepcopy/informer/lister/clientset
/home/cloud/ming/code-generator/generate-groups.sh \
all \
github.com/leemingeer/podgroup-scheduler-framework/pkg/generated \
github.com/leemingeer/podgroup-scheduler-framework/pkg/apis \
scheduling.ming.io:v1  \
--output-base ../../../ \
--go-header-file /home/cloud/ming/code-generator/hack/boilerplate.go.txt \
-v 10

# only generate crd
controller-gen  crd paths=./... output:crd:dir=config/crd