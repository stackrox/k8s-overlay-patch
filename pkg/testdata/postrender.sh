#!/usr/bin/env bash

set -o errexit

./k8s-overlay-patch -n default -p pkg/testdata/chart-patch.yaml