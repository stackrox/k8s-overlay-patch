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

package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	jsonpatch "github.com/evanphx/json-patch/v5" // nolint: staticcheck
	"github.com/kylelemons/godebug/diff"
	"github.com/stackrox/k8s-overlay-patch/pkg/protomarshal"
	"google.golang.org/protobuf/proto"
	yaml3 "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// ToYAML returns a YAML string representation of val, or the error string if an error occurs.
func ToYAML(val any) string {
	y, err := yaml.Marshal(val)
	if err != nil {
		return err.Error()
	}
	return string(y)
}

// UnmarshalWithJSONPB unmarshals y into out using gogo jsonpb (required for many proto defined structs).
func UnmarshalWithJSONPB(y string, out proto.Message, allowUnknownField bool) error {
	// Treat nothing as nothing.  If we called jsonpb.Unmarshaler it would return the same.
	if y == "" {
		return nil
	}
	jb, err := yaml.YAMLToJSON([]byte(y))
	if err != nil {
		return err
	}

	if allowUnknownField {
		err = protomarshal.UnmarshalAllowUnknown(jb, out)
	} else {
		err = protomarshal.Unmarshal(jb, out)
	}
	if err != nil {
		return err
	}
	return nil
}

// OverlayYAML patches the overlay tree over the base tree and returns the result. All trees are expressed as YAML
// strings.
func OverlayYAML(base, overlay string) (string, error) {
	if strings.TrimSpace(base) == "" {
		return overlay, nil
	}
	if strings.TrimSpace(overlay) == "" {
		return base, nil
	}
	bj, err := yaml.YAMLToJSON([]byte(base))
	if err != nil {
		return "", fmt.Errorf("yamlToJSON error in base: %s\n%s", err, bj)
	}
	oj, err := yaml.YAMLToJSON([]byte(overlay))
	if err != nil {
		return "", fmt.Errorf("yamlToJSON error in overlay: %s\n%s", err, oj)
	}
	if base == "" {
		bj = []byte("{}")
	}
	if overlay == "" {
		oj = []byte("{}")
	}

	merged, err := jsonpatch.MergePatch(bj, oj)
	if err != nil {
		return "", fmt.Errorf("json merge error (%s) for base object: \n%s\n override object: \n%s", err, bj, oj)
	}
	my, err := yaml.JSONToYAML(merged)
	if err != nil {
		return "", fmt.Errorf("jsonToYAML error (%s) for merged object: \n%s", err, merged)
	}

	return string(my), nil
}

// yamlDiff compares single YAML file
func yamlDiff(a, b string) string {
	ao, bo := make(map[string]any), make(map[string]any)
	if err := yaml.Unmarshal([]byte(a), &ao); err != nil {
		return err.Error()
	}
	if err := yaml.Unmarshal([]byte(b), &bo); err != nil {
		return err.Error()
	}

	ay, err := yaml.Marshal(ao)
	if err != nil {
		return err.Error()
	}
	by, err := yaml.Marshal(bo)
	if err != nil {
		return err.Error()
	}

	return diff.Diff(string(ay), string(by))
}

// yamlStringsToList yaml string parse to string list
func yamlStringsToList(str string) []string {
	reader := bufio.NewReader(strings.NewReader(str))
	decoder := yaml3.NewYAMLReader(reader)
	res := make([]string, 0)
	for {
		doc, err := decoder.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		chunk := bytes.TrimSpace(doc)
		res = append(res, string(chunk))
	}
	return res
}

// multiYamlDiffOutput multi yaml diff output format
func multiYamlDiffOutput(res, diff string) string {
	if res == "" {
		return diff
	}
	if diff == "" {
		return res
	}

	return res + "\n" + diff
}

func diffStringList(l1, l2 []string) string {
	var maxLen int
	var minLen int
	var l1Max bool
	res := ""
	if len(l1)-len(l2) > 0 {
		maxLen = len(l1)
		minLen = len(l2)
		l1Max = true
	} else {
		maxLen = len(l2)
		minLen = len(l1)
		l1Max = false
	}

	for i := 0; i < maxLen; i++ {
		d := ""
		if i >= minLen {
			if l1Max {
				d = yamlDiff(l1[i], "")
			} else {
				d = yamlDiff("", l2[i])
			}
		} else {
			d = yamlDiff(l1[i], l2[i])
		}
		res = multiYamlDiffOutput(res, d)
	}
	return res
}

// YAMLDiff compares multiple YAML files and single YAML file
func YAMLDiff(a, b string) string {
	al := yamlStringsToList(a)
	bl := yamlStringsToList(b)
	res := diffStringList(al, bl)

	return res
}

// IsYAMLEqual reports whether the YAML in strings a and b are equal.
func IsYAMLEqual(a, b string) bool {
	if strings.TrimSpace(a) == "" && strings.TrimSpace(b) == "" {
		return true
	}
	ajb, err := yaml.YAMLToJSON([]byte(a))
	if err != nil {
		//scope.Debugf("bad YAML in isYAMLEqual:\n%s", a)
		return false
	}
	bjb, err := yaml.YAMLToJSON([]byte(b))
	if err != nil {
		//scope.Debugf("bad YAML in isYAMLEqual:\n%s", b)
		return false
	}

	return bytes.Equal(ajb, bjb)
}

// IsYAMLEmpty reports whether the YAML string y is logically empty.
func IsYAMLEmpty(y string) bool {
	var yc []string
	for _, l := range strings.Split(y, "\n") {
		yt := strings.TrimSpace(l)
		if !strings.HasPrefix(yt, "#") && !strings.HasPrefix(yt, "---") {
			yc = append(yc, l)
		}
	}
	res := strings.TrimSpace(strings.Join(yc, "\n"))
	return res == "{}" || res == ""
}
