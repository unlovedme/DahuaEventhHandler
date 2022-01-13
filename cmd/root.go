
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

package cmd

import (
	"fmt"
	"os"

	"github.com/bitnami-labs/kubewatch/config"
	c "github.com/bitnami-labs/kubewatch/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kubewatch",
	Short: "A watcher for Kubernetes",
	Long: `
Kubewatch: A watcher for Kubernetes

kubewatch is a Kubernetes watcher that could publishes notification
to Slack/hipchat/mattermost/flock channels. It watches the cluster
for resource changes and notifies them through webhooks.

supported webhooks:
 - slack
 - hipchat
 - mattermost
 - flock
 - webhook
`,

	Run: func(cmd *cobra.Command, args []string) {
		config := &config.Config{}
		if err := config.Load(); err != nil {
			logrus.Fatal(err)