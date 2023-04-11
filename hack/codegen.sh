
#!/bin/bash

set -x

vendor/k8s.io/code-generator/generate-groups.sh all \
  github.com/obaydullahmhs/crd/pkg/client \
  github.com/obaydullahmhs/crd/pkg/apis \
  aadee.apps:v1alpha1 \
  --go-header-file /home/appscodepc/go/src/github.com/obaydullahmhs/crd/hack/boilerplate.go.txt
#
