// +k8s:deepcopy-gen=package

package v1alpha1

//go:generate deepcopy-gen -i github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -h hack/boilerplate.go.txt --trim-path-prefix github.com/stackrox/k8s-overlay-patch --output-base . --output-package github.com/stackrox/k8s-overlay-patch/pkg/apis/v1alpha1 -O zz_deepcopy
