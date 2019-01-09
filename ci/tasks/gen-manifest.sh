#!/usr/bin/env bash

go get sigs.k8s.io/kustomize
kustomize build ./examples/overlays > ../out/manifest.yaml
