#!/usr/bin/env bash

set -ex

go get sigs.k8s.io/kustomize
kustomize build repo/examples/overlays > out/manifest.yaml
