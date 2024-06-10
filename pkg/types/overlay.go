// Copyright 2019 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

type K8sObjectOverlay struct {
	// Resource API version.
	ApiVersion string `json:"apiVersion,omitempty"`
	// Resource kind.
	Kind string `json:"kind,omitempty"`
	// Name of resource.
	Name string `json:"name,omitempty"`
	// List of patches to apply to resource.
	Patches []*K8sObjectOverlayPatch `json:"patches,omitempty"`
	// Optional marks the overlay as optional. If the resource does not exist, the overlay is ignored.
	Optional bool `json:"optional,omitempty"`
}

type K8sObjectOverlayPatch struct {
	// Path of the form a.[key1:value1].b.[:value2]
	// Where [key1:value1] is a selector for a key-value pair to identify a list element and [:value] is a value
	// selector to identify a list element in a leaf list.
	// All path intermediate nodes must exist.
	Path string `json:"path,omitempty"`
	// Value to add, delete or replace.
	// For add, the path should be a new leaf.
	// For delete, value should be unset.
	// For replace, path should reference an existing node.
	// All values are strings but are converted into appropriate type based on schema.
	Value string `json:"value,omitempty"`
}

type OverlayObject struct {
	Overlays []*K8sObjectOverlay `json:"overlays,omitempty"`
}
