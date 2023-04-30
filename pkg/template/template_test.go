// Copyright © 2022 sealos.
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

package template

import (
	"bytes"
	"fmt"
	"github.com/labring/gh-rebot/pkg/types"
	"testing"
)

func TestTemplateSemverCompare(t *testing.T) {
	v, b, e := TryParse(`
scripts/changelog.sh {{.Repo.Fork}}
`)
	if e != nil {
		t.Errorf("parse err: %v", e)
	}
	if !b {
		t.Errorf("parse failed: %v", b)
	}
	types.GlobalsBotConfig = new(types.Config)
	types.GlobalsBotConfig.Repo.Fork = "cuisongliu/sealos"
	out := bytes.NewBuffer(nil)
	execErr := v.Execute(out, types.GlobalsBotConfig)
	if execErr != nil {
		t.Errorf("template exec err: %v", execErr)
	}

	fmt.Println(out)
}
