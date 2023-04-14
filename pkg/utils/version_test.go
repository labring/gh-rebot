/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"github.com/hashicorp/go-version"
	"testing"
)

func TestNewSemver(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"", true},
		{"1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"1.21.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"1.2.3-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", false},
		{"1.2.0.4-x.Y.0+metadata-width-hyphen", false},
		{"1.2.0-X-1.2.0+metadata~dist", false},
		{"1.2.3.4-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"v1.2.3", false},
		{"foo1.2.3", true},
		{"1.7rc2", true},
		{"v1.7rc2", true},
		{"1.0-", true},
	}

	for _, tc := range cases {
		_, err := version.NewSemver(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}
