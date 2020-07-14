# k8svent

Send Kubernetes pods as JSON to webhook endpoints. k8svent is run
from its Docker image in a Kubernetes cluster to send pod state
changes to the Atomist Kubernetes webhook endpoint for your Kubernetes
cluster integration.

## Running

At [https://go.atomist.com/](https://go.atomist.com/), create a
Kubernetes cluster integration and following the provided instructions
for deploying k8svent to your Kubernetes cluster.

## Webhook URLs

When running k8svent, webhook URLs can be specified in several ways:

-   The `--url` command-line option, which can be specified
    multiple times.

        $ k8svent --url=https://webhook.atomist.com/atomist/kube/teams/WORKSPACE_ID \
            --url=https://second.com/webhook

-   A comma-delimited list as the value of the `K8SVENT_WEBHOOKS`
    environment variable.

        $ K8SVENT_WEBHOOKS=https://webhook.atomist.com/0,https://webhook.atomist.com/1 k8svent

If webhooks are set using the `--url` command-line option, they
override any set by the `K8SVENT_WEBHOOKS` environment variable. In
other words, webhooks provided by the different methods are not
additive.

## Signing webhook payloads

k8svent can optionally sign the webhook payloads it sends using a
secret. The secret can be provided

-   The `--secret` command-line option.

        $ k8svent --secret=MyS3c43t

-   The value of the `K8SVENT_WEBHOOK_SECRET` environment variable.

        $ K8SVENT_WEBHOOK_SECRET=MyS3c43t k8svent

A secret provided on the command line takes precedence over one
provided via the environment variable. If a secret is provided, it is
used to sign the payloads send to all configured webhook endpoints.

## Webhook payload

k8svent sends payloads for _all_ pods, which it fetches using the
Kubernetes API. It periodically gets all pods and does its best to
send only the interesting ones, i.e., ones that have changed or that
are under duress. Each pod spec is serialized to JSON and sent to the
configured endpoint. The JSON structure is

```javascript
{
  "pod": {
    ... // k8s.io/api/core/v1.Pod
  }
}
```

The pod data structure is the same as you would see using `kubectl get pod POD -o json`.

## Updating

When running, k8svent periodically polls Docker Hub for tags and
checks the semantic version tags to see if any are newer than the
current running version. If the currently running version is a
release, it only checks tags that look like release versions. If the
currently running version is a prerelease, it checks all semantic
version tags for a newer version, which may be a release. If it
detects a newer version exists, it exits and lets Kubernetes pull the
new image and run it. To stay on the latest release, use the `latest`
tag. To use prerelease versions, use the `next` tag. To disable
updating, use a specific version tag.

## Developing

You can download, install, and develop locally using the normal Go
build tools.

```
$ go get github.com/atomist/k8svent
```

The source code will be under `$GOPATH/src/github.com/atomist/k8svent`.
If `$GOPATH/bin` is in your `PATH`, then the `k8svent` binary will be
in your path when the above command completes successfully. Then you
can run k8svent locally simply by invoking `k8svent` from your terminal.

If you make changes to the code, you can run tests using the Go
tooling

```
$ go test ./...
```

or you can use `make`

```
$ make test
```

To generate, build, test, install, vet, and lint, just run

```
$ make
```

---

Created by [Atomist][atomist].
Need Help? [Join our Slack team][slack].

[atomist]: https://atomist.com/ "Atomist - Automate All the Software Things"
[slack]: https://join.atomist.com/ "Atomist Community Slack Workspace"
