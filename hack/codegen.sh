
#!/bin/bash

set -x

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/obaydullahmhs/sample-controller/pkg/client \
  github.com/obaydullahmhs/sample-controller/pkg/apis \
  aadee.apps:v1alpha1 \
  --go-header-file /home/appscodepc/go/src/github.com/obaydullahmhs/sample-controller/hack/boilerplate.go.txt
#

#!/bin/bash

#set -x
#
#GROUP_NAME="aadee.apps"
#VERSION_NAME="v1alpha1"
#
#vendor/k8s.io/code-generator/generate-groups.sh all \
#github.com/obaydullahmhs/sample-controller/pkg/client \
#github.com/obaydullahmhs/sample-controller/pkg/apis \
#$GROUP_NAME:$VERSION_NAME \
#
#--go-header-file /home/appscodepc/go/src/github.com/obaydullahmhs/sample-controller/hack/boilerplate.go.txt