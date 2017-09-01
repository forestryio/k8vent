// Copyright © 2017 Atomist
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

	"github.com/atomisthq/k8vent/vent"
	"github.com/spf13/cobra"
)

const (
	defaultWebhookURL = "https://webhook.atomist.com/kube"
)

var (
	webhookURLs = []string{}
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the loop listening and emitting k8 events",
	Long: `Check for k8 events in an infinite loop and post them to
the configured webhooks.  You can provide the --url parameter multiple
times to send to multiple webhooks.

  $ k8vent --url=http://one.com/webhook --url=http://two.com/webhook

Alternatively, you can supply a comma-delimited list of webhook URLs in
the K8VENT_WEBHOOKS environment variable.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(webhookURLs) < 1 {
			webhookURLs = []string{defaultWebhookURL}
		}
		if err := vent.Vent(webhookURLs); err != nil {
			fmt.Fprintf(os.Stderr, "k8vent: venting failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().StringSliceVarP(&webhookURLs, "url", "u", []string{defaultWebhookURL}, "Send event to URL")
}
