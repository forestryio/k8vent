# k8vent

Send kubernetes pods as JSON to webhook endpoints.  k8vent is
typically run from its Docker image in a kubernetes cluster to send
pod state changes to the Atomist kubernetes webhook endpoint for your
team.

## Running

You can use [k8vent-deployment.json][k8vent-deployment] as a starting
point for running k8vent in your kubernetes cluster.  Be sure to
update the value of the `K8VENT_WEBHOOKS` environment variable,
replacing the last element of the URL with the ID of your Slack team.
Slack team IDs start with a `T` and are nine characters long,
consisting of digits and capital letters.  You can get your team ID
from https://app.atomist.com/teams or by sending `team` as a message
to the Atomist bot, e.g., `@atomist team`, in Slack.

You can optionally change the value of the
`ATOMIST_ENVIRONMENT` environment variable to a meaningful name for
your kubernetes cluster, e.g., "prod" or "staging".

Once you have made the changes, you can create the k8vent deployment
as you normally would.

    $ kubectl create -f k8vent-deployment.json

If you have [jq][] installed, you can update to a new version of
k8vent with the following command

    $ kubectl get deployment k8vent -o json \
        | jq ".spec.template.spec.containers[0].image=\"atomist/k8vent:M.N.P\"" \
        | kubectl replace -f -

replacing `M.N.P` with the [latest version of k8vent][latest].

[k8vent-deployment]: k8vent-deployment.json (k8vent Kubernetes Deployment Spec)
[jq]: https://stedolan.github.io/jq/ (jq)
[latest]: https://github.com/atomist/k8vent/releases/latest (k8vent Current Release)

## Webhook payload

k8vent subscribes to all pods via the kubernetes watch API.  When a
pod changes, e.g., it is created, deleted or otherwise changes state,
the current pod structure and the environment of the `k8vent` pod is
serialized as JSON and sent to the configured webhook endpoints.  The
JSON structure is

```javascript
{
  "pod": {
    ... // k8s.io/api/core/v1.Pod
  },
  "env": {
    "ENV_VAR1": "value1",
    "ENV_VAR2": "value2",
    ...
  }
}
```

The pod data strcuture is the same as you would see using `kubectl get
pod POD -o json`.

## Developing

You can download, install, and develop locally using the normal Go
build tools.

    $ go get github.com/atomist/k8vent

The source code will be under `$GOPATH/src/github.com/atomist/k8vent`.
If `$GOPATH/bin` is in your `PATH`, then the `k8vent` binary will be
in your path when the above command completes successfully.  Then you
can run k8vent locally simply by invoking `k8vent` from your terminal.

When running from the command line, webhook URLs can be specified in
several ways:

-   Using the `--url` command-line option, which can be specified
    multiple times.

        $ k8vent --url=http://first.com/webhook --url=https://second.com/webhook

-   A comma-delimited list as the value of the `K8VENT_WEBHOOKS`
    environment variable.

        $ K8VENT_WEBHOOKS=http://first.com/webhook,https://second.com/webhook k8vent

If webhooks are set using the `--url` command-line option, they
override any set by the `K8VENT_WEBHOOKS` environment variable.  If no
webhook URLs are provided, `k8vent` exits with an error.  In other
words, webhooks provided by the different methods are not additive.

If you make changes to the code, you can run tests using the Go
tooling

    $ go test ./...

or you can use `make`

    $ make

---

Created by [Atomist][atomist].
Need Help?  [Join our Slack team][slack].

[atomist]: https://www.atomist.com/
[slack]: https://join.atomist.com/
