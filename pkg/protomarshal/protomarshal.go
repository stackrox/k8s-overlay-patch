// Copyright Istio Authors
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

// Package protomarshal provides operations to marshal and unmarshal protobuf objects.
// Unlike the rest of this repo, which uses the new google.golang.org/protobuf API, this package
// explicitly uses the legacy jsonpb package. This is due to a number of compatibility concerns with the new API:
// * https://github.com/golang/protobuf/issues/1374
// * https://github.com/golang/protobuf/issues/1373
package protomarshal

import (
	jsonpb "google.golang.org/protobuf/encoding/protojson" // nolint: depguard
	"google.golang.org/protobuf/proto"
)

func Unmarshal(b []byte, m proto.Message) error {
	return jsonpb.UnmarshalOptions{DiscardUnknown: false}.Unmarshal(b, m)
}

func UnmarshalAllowUnknown(b []byte, m proto.Message) error {
	return jsonpb.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(b, m)
}
