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

package config

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"strings"
)

func ParseConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := yaml.NewYAMLOrJSONDecoder(file, 1024)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	if token, ok := os.LookupEnv("GH_TOKEN"); ok {
		config.SetToken(token)
	}
	if config.Repo.Org {
		config.Repo.OrgCommand = fmt.Sprintf(" --org  %s ", strings.SplitN(config.GetRepoName(), "/", 2)[0])
	}
	return config, nil
}

func LoadConfig(cfg string) (*Config, error) {
	configPaths := []string{".github/gh-bot.yml", ".github/gh-bot.yaml", cfg}
	for _, configPath := range configPaths {
		if _, err := os.Stat(configPath); err == nil {
			config, err := ParseConfig(configPath)
			if err != nil {
				return nil, fmt.Errorf("error parsing config file %s: %v", configPath, err)
			}
			return config, nil
		}
	}
	return nil, fmt.Errorf("no valid config file found")
}
