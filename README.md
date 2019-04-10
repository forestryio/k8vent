# k8vent

[![atomist sdm goals](https://badge.atomist.com/T29E48P34/atomist/k8vent/f4268ca9-b9f2-4717-902b-75892c2ec9ed)](https://app.atomist.com/workspace/T29E48P34)

Send kubernetes pods as JSON to webhook endpoints.  k8vent is
typically run from its Docker image in a kubernetes cluster to send
pod state changes to the Atomist kubernetes webhook endpoint for your
Atomist workspace.

## Running

See the [Atomist Kubernetes documentation][atomist-kube] for detailed
instructions on using Atomist with Kubernetes.  Briefly, if you
already have an [Atomist workspace][atomist-getting-started], you can
run the following commands to create the necessary resources in your
Kubernetes cluster.  Replace `CLUSTER_ENV` with a unique name for your
Kubernetes cluster and `WORKSPACE_ID` with your Atomist workspace ID.

[atomist-kube]: https://docs.atomist.com/pack/kubernetes/ (Atomist - Kubernetes)
[atomist-getting-started]: https://docs.atomist.com/user/ (Atomist - Getting Started)

```
$ kubectl apply -f https://raw.githubusercontent.com/atomist/k8vent/master/kube/kubectl/cluster-wide.yaml
$ kubectl create secret --namespace=k8vent generic k8vent --from-literal=environment=CLUSTER_ENV \
    --from-literal=webhooks=https://webhook.atomist.com/atomist/kube/teams/WORKSPACE_ID
```

## Webhook URLs

When running k8vent, webhook URLs can be specified in several ways:

-   Pod-specific webhook URLs in the pod's metadata annotations.  Use
    "atomist.com/k8vent" as the annotation key and the value should be
    a properly escaped JSON object with the key "webhooks" whose value
    is an array of webhook URLs.  For example:

        metadata:
          annotations:
            atomist.com/k8vent: '{"webhooks":["https://webhook.atomist.com/atomist/kube/teams/WORKSPACE_ID"]}'

-   The `--url` command-line option, which can be specified
    multiple times.

        $ k8vent --url=https://webhook.atomist.com/atomist/kube/teams/WORKSPACE_ID \
            --url=https://second.com/webhook

-   A comma-delimited list as the value of the `K8VENT_WEBHOOKS`
    environment variable.

        $ K8VENT_WEBHOOKS=https://webhook.atomist.com/atomist/kube/teams/WORKSPACE_ID,https://second.com/webhook k8vent

If webhooks are provided in the pod spec, they override any provided
on the command line or by the environment.  If webhooks are set using
the `--url` command-line option, they override any set by the
`K8VENT_WEBHOOKS` environment variable.  In other words, webhooks
provided by the different methods are not additive.

## Environment

You are able to tag each pod event sent with an _environment_ string,
e.g., "production", "qa", or "testing".  This value is used in Atomist
lifecycle messages to differentiate between pods deployed in different
clusters, namespaces, etc.  You can provide this value two different
ways:

-   Pod-specific environment in the pod's metadata annotations.  Use
    "atomist.com/k8vent" as the annotation key and the value should be
    a properly escaped JSON object with the key "environment" whose
    value is the environment string.  For example:

        metadata:
          annotations:
            atomist.com/k8vent: '{"environment":"production"}'

-   The environment variable `ATOMIST_ENVIRONMENT`

        $ ATOMIST_ENVIRONMENT=production k8vent

The pod-specific annotation overrides any value of
`ATOMIST_ENVIRONMENT` in the k8vent process' environment.

## Linking to Atomist lifecycle events

To get Kubernetes pod events to display as part of the normal Git
activity messages you see from Atomist in Slack channels or the
Atomist dashboard event stream, you must tell Atomist what Docker
images are connected with what commits.  The link between a commit and
a Docker image is created by POSTing data to the Atomist webhook
endpoint
`https://webhook.atomist.com/atomist/link-image/teams/WORKSPACE_ID`,
where `WORKSPACE_ID` should be replaced with the same Atomist
workspace ID used in the `K8VENT_WEBHOOKS` environment variable above.
The POST data should be JSON of the form:

```json
{
  "git": {
    "owner": "REPO_OWNER",
    "repo": "REPO_NAME",
    "sha": "COMMIT_SHA"
  },
  "docker": {
    "image":"DOCKER_IMAGE_TAG"
  },
  "type":"link-image"
}
```

where `REPO_OWNER` is the repository owner, i.e., the user or
organization, `REPO_NAME` is the name of the repository, `COMMIT_SHA`
is the full SHA of the commit from which the Docker image was created,
and `DOCKER_IMAGE_TAG` is the full tag for the Docker image, i.e.,
OWNER/IMAGE:VERSION for Docker Hub images or
REGISTRY/OWNER/IMAGE:VERSION for images in other Docker registries.

If you have a shell script executing your CI build that creates your
Docker image, you can add the following command after the Docker image
has been pushed to the registry.

```
$ curl -s -f -X POST -H "Content-Type: application/json" \
    --data-binary '{"git":{…},"docker":{…},"type":"link-image"}' \
    https://webhook.atomist.com/atomist/link-image/teams/WORKSPACE_ID
```

replacing the ellipses with the appropriate JSON and `WORKSPACE_ID` with
the appropriate Atomist workspace ID.

## Webhook payload

k8vent subscribes to _all_ pods via the kubernetes watch API.  When a
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

```
$ go get github.com/atomist/k8vent
```

The source code will be under `$GOPATH/src/github.com/atomist/k8vent`.
If `$GOPATH/bin` is in your `PATH`, then the `k8vent` binary will be
in your path when the above command completes successfully.  Then you
can run k8vent locally simply by invoking `k8vent` from your terminal.

If you make changes to the code, you can run tests using the Go
tooling

```
$ go test ./...
```

or you can use `make`

```
$ make test
```

To build, test, install, and vet, just run

```
$ make
```

---

Created by [Atomist][atomist].
Need Help?  [Join our Slack team][slack].

[atomist]: https://atomist.com/ (Atomist - How Teams Deliver Software)
[slack]: https://join.atomist.com/ (Atomist Community Slack Workspace)
