# k8s-overlay-patch
Kubernetes overlay patches

## Disclaimer
This codebase was mostly plucked from https://github.com/istio/istio

## Purpose
Users often need to customize the resources generated from a kubernetes operator/Custom Resource. Though, for every new
configuration option that needs to be added, new fields need to be added to the Custom Resource. 
This is time-consuming and can lead to a bloated Custom Resource. 

This library implements logic so that these overlays are a first-class citizen in a CRD. This allows
users to customize the output of a kubernetes operator by overlaying patches on top of the generated output, without
having to modify the CRD.

## Command-line usage

```
Usage:
  k8s-overlay-patch [flags]

Flags:
  -h, --help                   help for k8s-overlay-patch
  -m, --manifest-file string   File containing the rendered manifests to patch
  -n, --namespace string       Namespace to use when patching the manifests
  -o, --out string             File to write the patched manifests to
  -p, --patch-file string      File containing the patch to apply
```


## Usage as helm post-renderer

Example

```
helm template test pkg/testdata/chart --debug --dry-run --post-renderer pkg/testdata/postrender.sh
```

```
# build the binary and place it in your path
go build

# Create a postrender.sh file. Helm doesn't allow passing
# arguments to postrenderers.
#!/bin/bash
k8s-overlay-patch -n default -p patch.yaml

# Make the file executable
chmod +x postrender.sh

# Call helm with the postrender.sh file
helm template [NAME] [CHART] --post-renderer postrender.sh

```

## Example

#### Adding the overlay patch to the CRD
```go
package v1alpha1

type MySpec struct {
	// Overlays is the list of overlay patches to apply to resource.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Patches",order=4
	Overlays []K8sObjectOverlay `json:"overlays,omitempty"`
}


type K8sObjectOverlay struct {
	// Resource API version.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="API Version",order=1
	ApiVersion string `json:"apiVersion,omitempty"`
	// Resource kind.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Kind",order=2
	Kind string `json:"kind,omitempty"`
	// Name of resource.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name",order=3
	Name string `json:"name,omitempty"`
	// List of patches to apply to resource.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Patches",order=4
	Patches []*K8sObjectOverlayPatch `json:"patches,omitempty"`
}

type K8sObjectOverlayPatch struct {
	// Path of the form a.[key1:value1].b.[:value2]
	// Where [key1:value1] is a selector for a key-value pair to identify a list element and [:value] is a value
	// selector to identify a list element in a leaf list.
	// All path intermediate nodes must exist.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Path",order=1
	Path string `json:"path,omitempty"`
	// Value to add, delete or replace.
	// For add, the path should be a new leaf.
	// For delete, value should be unset.
	// For replace, path should reference an existing node.
	// All values are strings but are converted into appropriate type based on schema.
	//+operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Value",order=2
	Value string `json:"value,omitempty"`
}

```


#### Using the overlay patch in the reconciler
```go
package reconciler

import (
	"github.com/stackrox/k8s-overlay-patch/pkg/patch"
	"github.com/stackrox/k8s-overlay-patch/pkg/types"
)

func (r reconciler) Reconcile(obj v1alpha1.MyObject){
	manifests := renderManifests(obj)
	patched := patch.YAMLManifestPatch(manifests, obj.Namespace, mapOverlays(obj.Spec.Overlays))
	// ...
}

func mapOverlays(overlays []*v1alpha1.K8sObjectOverlay) []*types.K8sObjectOverlay {
	out := make([]*types.K8sObjectOverlay, len(overlays))
	for i, o := range overlays {
		out[i] = &types.K8sObjectOverlay{
			ApiVersion: o.ApiVersion,
			Kind:       o.Kind,
			Name:       o.Name,
			Patches:    mapOverlayPatches(o.Patches),
		}
	}
	return out
}

func mapOverlayPatches(patches []*v1alpha1.K8sObjectOverlayPatch) []*types.K8sObjectOverlayPatch {
	out := make([]*types.K8sObjectOverlayPatch, len(patches))
	for i, p := range patches {
		out[i] = &types.K8sObjectOverlayPatch{
			Path:  p.Path,
			Value: p.Value,
		}
	}
	return out
}


```

#### Example CRD
```yaml
apiVersion: blah.com/v1alpha1
kind: MyObject
metadata:
  name: my-object
spec:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi        
  overlays:
  - apiVersion: v1
    kind: ConfigMap
    name: my-config-map
    patches:
    - path: data.foo
      value: bar
    - path: data.baz
      value: qux
```