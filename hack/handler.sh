#!/bin/sh

# generate deepcopy/informer/lister/clientset
/home/cloud/ming/code-generator/generate-groups.sh \
all \
github.com/leemingeer/podgroup-scheduler-framework/pkg/generated \
github.com/leemingeer/podgroup-scheduler-framework/pkg/apis \
scheduling:v1  \
--output-base ../../../ \
--go-header-file /home/cloud/ming/code-generator/hack/boilerplate.go.txt \
-v 10

# for internal type
/home/cloud/ming/code-generator/generate-internal-groups.sh \
"deepcopy,defaulter,conversion" \
github.com/leemingeer/podgroup-scheduler-framework/pkg/generated \
github.com/leemingeer/podgroup-scheduler-framework/pkg/apis \
github.com/leemingeer/podgroup-scheduler-framework/pkg/apis \
"config:v1beta1"  \
--output-base ../../../ \
--go-header-file /home/cloud/ming/code-generator/hack/boilerplate.go.txt \
-v 10


# only generate crd
/root/go/bin/controller-gen  crd paths=./... output:crd:dir=config/crd

# only generate config api deepcopy
/root/go/bin/controller-gen object paths=./pkg/apis/config/v1beta1/types.go
