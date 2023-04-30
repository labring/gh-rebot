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

package types

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/labring/gh-rebot/pkg/utils"
	"time"
)

type RetryShell string
type RetrySecretShell string
type SecretShell string

func ExecShellForAny(secrets ...string) func(shells []any) error {
	return func(shells []any) error {
		// 设置重试策略
		exponentialBackoff := backoff.NewExponentialBackOff()
		exponentialBackoff.MaxElapsedTime = 15 * time.Second
		for _, sh := range shells {
			if s, ok := sh.(RetryShell); ok {
				if err := backoff.Retry(func() error {
					return utils.RunCommand("bash", "-c", string(s))
				}, exponentialBackoff); err != nil {
					return err
				}
			}
			if s, ok := sh.(RetrySecretShell); ok {
				if err := backoff.Retry(func() error {
					return utils.RunCommandInSecret(string(s), secrets)
				}, exponentialBackoff); err != nil {
					return err
				}
			}
			if s, ok := sh.(SecretShell); ok {
				if err := utils.RunCommandInSecret(string(s), secrets); err != nil {
					return err
				}
			}
			if s, ok := sh.(string); ok {
				if err := utils.RunCommand("bash", "-c", s); err != nil {
					return err
				}
			}
		}
		return nil
	}
}
