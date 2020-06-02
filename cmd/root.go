// Copyright © 2020 Atomist
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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/atomist/k8svent/vent"
)

var cfgFile string

var (
	namespace     string
	webhookSecret string
	webhookURLs   = []string{}
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "k8svent",
	Short: "Send kubernetes pod state changes to webhook",
	Long: `Watch for kubernetes pod state changes and post them to the configured
webhooks.

You can provide the --url parameter multiple times to send to multiple
webhooks.

  $ k8svent --url=http://one.com/webhook --url=http://two.com/webhook

Alternatively, you can supply a comma-delimited list of webhook URLs
in the K8SVENT_WEBHOOKS environment variable or provide them in the pod
annotations.

By default k8svent watches pods in all namespaces.  If the --namespace
or K8SVENT_NAMESPACE environment variable is provided, only pods in
that namespace are reported on.

By default k8svent does not sign the webhook payloads.  If the
--secret or K8SVENT_WEBHOOK_SECRET environment variable is provided,
webhook payloads are signed using HMAC/SHA-1.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := vent.Vent(webhookURLs, namespace, webhookSecret); err != nil {
			fmt.Fprintf(os.Stderr, "k8svent: venting failed: %v\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	//RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k8svent.yaml)")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Only watch pods in NAMESPACE")
	RootCmd.PersistentFlags().StringVarP(&webhookSecret, "secret", "s", "", "Sign webhook payloads using SECRET")
	RootCmd.PersistentFlags().StringSliceVarP(&webhookURLs, "url", "u", []string{}, "Send event to URL")
}

const webhookEnv = "K8SVENT_WEBHOOKS"
const namespaceEnv = "K8SVENT_NAMESPACE"
const webhookSecretEnv = "K8SVENT_WEBHOOK_SECRET"

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if os.Getenv(webhookEnv) != "" {
		webhookURLs = strings.Split(os.Getenv(webhookEnv), ",")
	}
	if os.Getenv(namespaceEnv) != "" {
		namespace = os.Getenv(namespaceEnv)
	}
	if os.Getenv(webhookSecretEnv) != "" {
		webhookSecret = os.Getenv(webhookSecretEnv)
	}

	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".k8svent") // name of config file (without extension)
	viper.AddConfigPath("$HOME")    // adding home directory as first search path
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
