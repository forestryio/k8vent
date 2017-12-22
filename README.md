# k8vent

Send kuberenetes events to a webhook.

## Installing

You can install using the normal Go method.

    $ go get github.com/atomist/k8vent

If `$GO_HOME/bin` is in your `PATH`, then the `k8vent` binary will be
in your path when the above command completes successfully.

## Running

The program periodically polls the kubernetes Event API and sends new
events and events whose Count has increased to webhooks.  It should be
run inside a kubernetes cluster.

Webhook URLs can be specified in several ways:

-   Using the `--url` command-line option, which can be specified
    multiple times.

        $ k8vent --url=http://first.com/webhook --url=https://second.com/webhook

-   A comma-delimited list as the value of the `K8VENT_WEBHOOKS`
    environment variable.

        $ K8VENT_WEBHOOKS=http://first.com/webhook,https://second.com/webhook k8vent

-   Accept the default webhook: https://webhook.atomist.com/kube

If webhooks are set using the `--url` command-line option, they
override any set by the `K8VENT_WEBHOOKS` environment variable.  If
any are set in the environment variable, they override the default
value.  In other words, webhooks provided by the different methods are
not additive.

---

Created by [Atomist][atomist].
Need Help?  [Join our Slack team][slack].

[atomist]: https://www.atomist.com/
[slack]: https://join.atomist.com/
