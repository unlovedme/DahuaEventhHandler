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
	"github.com/bitnami-labs/kubewatch/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// resourceConfigCmd represents the resource subcommand
var resourceConfigCmd = &cobra.Command{
	Use:   "resource",
	Short: "manage resources to be watched",
	Long: `
manage resources to be watched`,
	Run: func(cmd *cobra.Command, args []string) {

		// warn for too few arguments
		if len(args) < 2 {
			logrus.Warn("Too few arguments to Command \"resource\".\nMinimum 2 arguments required: subcommand, resource flags")
		}
		// display help
		cmd.Help()
	},
}

// resourceConfigAddCmd represents the resource add subcommand
var resourceConfigAddCmd = &cobra.Command{
	Use:   "add",
	Short: "adds specific resources to be watched",
	Long: `
adds specific resources to be watched`,
	Run: func(cmd *cobra.Command, args []string) {
		conf, err := config.New()
		if err != nil {
			logrus.Fatal(err)
		}

		// add resource to config
		configureResource("add", cmd, conf)
	},
}

// resourceConfigRemoveCmd represents the resource remove subcommand
var resourceConfigRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove specific resources being watched",
	Long: `
remove 