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
	Node                  bool `json:"node"`
	ClusterRole           bool `json:"clusterrole"`
	ServiceAccount        bool `json:"sa"`
	PersistentVolume      bool `json:"pv"`
	Namespace             bool `json:"ns"`
	Secret                bool `json:"secret"`
	ConfigMap             bool `json:"configmap"`
	Ingress               bool `json:"ing"`
}

// Config struct contains kubewatch configuration
type Config struct {
	// Handlers know how to send notifications to specific services.
	Handler Handler `json:"handler"`

	//Reason   []string `json:"reason"`

	// Resources to watch.
	Resource Resource `json:"resource"`

	// For watching specific namespace, leave it empty for watching all.
	// this config is ignored when watching namespaces
	Namespace string `json:"namespace,omitempty"`
}

// Slack contains slack configuration
type Slack struct {
	// Slack "legacy" API token.
	Token string `json:"token"`
	// Slack channel.
	Channel string `json:"channel"`
	// Title of the message.
	Title string `json:"title"`
}

// Hipchat contains hipchat configuration
type Hipchat struct {
	// Hipchat token.
	Token string `json:"token"`
	// Room name.
	Room string `json:"room"`
	// URL of the hipchat server.
	Url string `json:"url"`
}

// Mattermost contains mattermost configuration
type Mattermost struct {
	Channel  string `json:"room"`
	Url      string `json:"url"`
	Username string `json:"username"`
}

// Flock contains flock configuration
type Flock struct {
	// URL of the flock API.
	Url string `json:"url"`
}

// Webhook contains webhook configuration
type Webhook struct {
	// Webhook URL.
	Url string `json:"url"`
}

// MSTeams contains MSTeams configuration
type MSTeams struct {
	// MSTeams API Webhook URL.
	WebhookURL string `json:"webhookurl"`
}

// SMTP contains SMTP configuration.
type SMTP struct {
	// Destination e-mail address.
	To string `json:"to" yaml:"to,omitempty"`
	// Sender e-mail address .
	From string `json:"from" yaml:"from,omitempty"`
	// Smarthost, aka "SMTP server"; address of server used to send email.
	Smarthost string `json:"smarthost" yaml:"smarthost,omitempty"`
	// Subject of the outgoing emails.
	Subject string `json:"subject" yaml:"subject,omitempty"`
	// Extra e-mail headers to be added to all outgoing messages.
	Headers map[string]string `json:"headers" yaml:"hea