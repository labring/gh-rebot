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

var (
	GlobalsConfig *Config
)

type Bot struct {
	AllowOps []string `json:"allowOps"`
	Email    string   `json:"email"`
	Username string   `json:"username"`
}

// GetAllowOps returns the triggers for the bot
func (r *Config) GetAllowOps() []string {
	return r.Bot.AllowOps
}

// GetEmail returns the email for the bot
func (r *Config) GetEmail() string {
	return r.Bot.Email
}

// GetUsername returns the username for the bot
func (r *Config) GetUsername() string {
	return r.Bot.Username
}

type Repo struct {
	Org        bool   `json:"org"`
	OrgCommand string `json:"-"`
	Name       string `json:"name"`
	Fork       string `json:"fork"`
}

// GetOrgCommand returns the org command for the repo
func (r *Config) GetOrgCommand() string {
	return r.Repo.OrgCommand
}

// GetRepoName returns the name for the repo
func (r *Config) GetRepoName() string {
	return r.Repo.Name
}

// GetForkName returns the fork for the repo
func (r *Config) GetForkName() string {
	return r.Repo.Fork
}

type Config struct {
	Version string            `json:"version"`
	Debug   bool              `json:"debug"`
	Bot     Bot               `json:"bot"`
	Repo    Repo              `json:"repo"`
	Message map[string]string `json:"message"`
	Token   string            `json:"-"`
}

// GetDebug returns the debug for the config
func (c *Config) GetDebug() bool {
	return c.Debug
}

// GetToken returns the token for the config
func (c *Config) GetToken() string {
	return c.Token
}

// SetToken sets the token for the config
func (c *Config) SetToken(token string) {
	c.Token = token
}

// GetMessages returns the messages for the config
func (c *Config) GetMessages() map[string]string {
	return c.Message
}

func (c *Config) GetMessage(key, defaultVal string) string {
	if c.Message[key] != "" {
		return c.Message[key]
	}
	return defaultVal
}
