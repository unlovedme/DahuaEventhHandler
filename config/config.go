/*
Copyright 2016 Skippbox, Ltd.

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

//go:generate bash -c "go install ../tools/yannotated && yannotated -o sample.go -format go -package config -type Config"

package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

var (
	// ConfigFileName stores file of config
	ConfigFileName = ".kubewatch.yaml"

	// ConfigSample is a sample configuration file.
	ConfigSample = yannotated
)

// Handler contains handler configuration
type Handler struct {
	Slack      Slack      `json:"slack"`
	Hipchat    Hipchat    `json:"hipchat"`
	Mattermost Mattermost `json:"mattermost"`
	Flock      Flock      `json:"flock"`
	Webhook    Webhook    `json:"webhook"`
	MSTeams    MSTeams    `json:"msteams"`
	SMTP       SMTP       `json:"smtp"`
}

// Resource contains resource configuration
type Resource struct {
	Deployment            bool `json:"deployment"`
	ReplicationController bool `json:"rc"`
	ReplicaSet            bool `json:"rs"`
	DaemonSet             bool `json:"ds"`
	Services              bool `json:"svc"`
	Pod                   bool `json:"po"`
	Job                   bool `json:"job"`
	Node                  bool `json:"node