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
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/cuisongliu/logger"
	"os"
	"os/exec"
	"strings"
	"time"
)

type RetryShell string
type RetrySecretShell string
type SecretShell string
type SkipShell string

func ExecShellForAny(secrets ...string) func(shells []any) error {
	return func(shells []any) error {
		// 设置重试策略
		exponentialBackoff := backoff.NewExponentialBackOff()
		exponentialBackoff.MaxElapsedTime = 15 * time.Second
		for _, sh := range shells {
			if s, ok := sh.(RetryShell); ok {
				if err := backoff.Retry(func() error {
					return RunCommand("bash", "-c", string(s))
				}, exponentialBackoff); err != nil {
					return err
				}
			}
			if s, ok := sh.(RetrySecretShell); ok {
				if err := backoff.Retry(func() error {
					return RunCommandInSecret(string(s), secrets)
				}, exponentialBackoff); err != nil {
					return err
				}
			}
			if s, ok := sh.(SecretShell); ok {
				if err := RunCommandInSecret(string(s), secrets); err != nil {
					return err
				}
			}
			if _, ok := sh.(SkipShell); ok {
				continue
			}
			if s, ok := sh.(string); ok {
				if err := RunCommand("bash", "-c", s); err != nil {
					return err
				}
			}

		}
		return nil
	}
}

func RunCommand(command string, args ...string) error {
	logger.Debug("Running command: %s %s", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed (%s): %v", strings.Join(cmd.Args, " "), err)
	}
	return nil
}

func RunCommandInSecret(command string, secrets []string) error {
	var disPlayCommand string
	for _, secret := range secrets {
		disPlayCommand = strings.ReplaceAll(command, secret, "******")
	}
	logger.Debug("Running command: %s", disPlayCommand)
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed (%s): %v", disPlayCommand, err)
	}
	return nil
}

func RunCommandWithOutput(cmd string, removeLine bool) (string, error) {
	logger.Debug("Running command with output: %s", cmd)
	// nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	result, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput() // #nosec
	if removeLine {
		out := getOnelineResult(string(result), "")
		logger.Debug("Command output: %s", out)
		return out, err
	}
	out := string(result)
	logger.Debug("Command output: %s", out)
	return out, err
}

func getOnelineResult(output string, sep string) string {
	return strings.ReplaceAll(strings.ReplaceAll(output, "\r\n", sep), "\n", sep)
}
