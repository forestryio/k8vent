# k8vent

[![Build Status](https://travis-ci.org/atomist/k8vent.svg?branch=master)](https://travis-ci.org/atomist/k8vent)

Send kubernetes pods as JSON to webhook endpoints.  k8vent is
typically run from its Docker image in a kubernetes cluster to send
pod state changes to the Atomist kubernetes webhook endpoint for your
team.

## Running

You can use the Kubernetes resource files in the [kube
directory][kube] as a starting point for deploying k8vent in your
kubernetes cluster.

k8vent needs read access to pods to operate normally.  It uses the
Kubernetes "in-cluster client" to authenticate against the Kubernetes
API.  Depending on whether your cluster is using [role-based access
control (RBAC)][rbac].

Before deploying either with or without RBAC, you will need your
Atomist workspace/team ID.  You can get your Atomist team/workspace ID
from its settings page on https://app.atomist.com/ or by sending
`team` as a message to the Atomist bot, e.g., `@atomist team`, in
Slack.

You can  optionally provide  an environment  name for  your kubernetes
cluster, e.g., "production", "qa", or "testing".

[kube]: ./kube/ (k8vent Kubernetes Resources)
[rbac]: https://kubernetes.io/docs/admin/authorization/rbac/ (Kubernetes RBAC)
[gke-rbac]: https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control (GKE RBAC)

You can deploy k8vent with the following commands

```console
$ kubectl apply -f https://raw.githubusercontent.com/atomist/k8vent/master/kube/kubectl/cluster-wide.yaml
$ kubectl create secret --namespace=k8vent generic k8vent --from-literal=environment=CLUSTER_ENV \
    --from-literal=webhooks=https://webhook.atomist.com/atomist/kube/teams/TEAM_ID
```

replacing `TEAM_ID` with your Atomist team ID and `CLUSTER_ENV` with
your desired environment name.

If you get the following error when running the first command,

```
Error from server (Forbidden): error when creating "rbac.yaml": clusterroles.rbac.authorization.k8s.io "k8vent-clusterrole" is forbidden: attempt to grant extra privileges: [PolicyRule{Resources:["pods"], APIGroups:[""], Verbs:["get"]} PolicyRule{Resources:["pods"], APIGroups:[""], Verbs:["list"]} PolicyRule{Resources:["pods"], APIGroups:[""], Verbs:["watch"]}] user=&{YOUR_USER  [system:authenticated] map[]} ownerrules=[PolicyRule{Resources:["selfsubjectaccessreviews"], APIGroups:["authorization.k8s.io"], Verbs:["create"]} PolicyRule{NonResourceURLs:["/api" "/api/*" "/apis" "/apis/*" "/healthz" "/swagger-2.0.0.pb-v1" "/swagger.json" "/swaggerapi" "/swaggerapi/*" "/version"], Verbs:["get"]}] ruleResolutionErrors=[]
```

then your Kubernetes user does not have cluster-admin role privileges
on your cluster.  You will either need to ask someone who has
cluster-admin privileges on the cluster to create the RBAC resources
or try to escalate your privileges with the following command.

```console
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin \
    --user YOUR_USER
```

If you are running on GKE, you can supply your user name using the
`gcloud` utility.

```console
$ kubectl create clusterrolebinding cluster-admin-binding --clusterrole cluster-admin \
    --user $(gcloud config get-value account)
```

Then run the commands again.

### Updating

You can update to a new version of k8vent with the following command.

```console
$ kubectl patch --namespace=k8vent deployment k8vent \
    --patch '{"spec":{"template":{"spec":{"containers":[{"name":"k8vent","image":"atomist/k8vent:M.N.P"}]}}}}'
```

replacing `M.N.P` with the [latest version of k8vent][latest].

[latest]: https://github.com/atomist/k8vent/releases/latest (k8vent Current Release)

## Webhook URLs

When running k8vent, webhook URLs can be specified in several ways:

-   Pod-specific webhook URLs in the pod's metadata annotations.  Use
    "atomist.com/k8vent" as the annotation key and the value should be
    a properly escaped JSON object with the key "webhooks" whose value
    is an array of webhook URLs.  For example:

        "metadata": {
          "labels": {
            "app": "my-app",
          },
          "annotations": {
            "atomist.com/k8vent": "{\"webhooks\":[\"https://webhook.atomist.com/atomist/kube/teams/TEAM_ID\"]}"
          }
        },

-   The `--url` command-line option, which can be specified
    multiple times.

        $ k8vent --url=https://webhook.atomist.com/atomist/kube/teams/TEAM_ID \
            --url=https://second.com/webhook

-   A comma-delimited list as the value of the `K8VENT_WEBHOOKS`
    environment variable.

        $ K8VENT_WEBHOOKS=https://webhook.atomist.com/atomist/kube/teams/TEAM_ID,https://second.com/webhook k8vent

If webhooks are provided in the pod spec, they override any provided
on the command line or by the environment.  If webhooks are set using
the `--url` command-line option, they override any set by the
`K8VENT_WEBHOOKS` environment variable.  In other words, webhooks
provided by the different methods are not additive.  If no webhook
URLs are provided, `k8vent` exits with an error.

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

        "metadata": {
          "labels": {
            "app": "my-app",
          },
          "annotations": {
            "atomist.com/k8vent": "{\"environment\":\"production\"}"
          }
        },

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
`https://webhook.atomist.com/atomist/link-image/teams/TEAM_ID`, where
`TEAM_ID` should be replaced with the same team ID used in the
`K8VENT_WEBHOOKS` environment variable above.  The POST data should be
JSON of the form:

```javascript
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

```shell
curl -s -f -X POST -H "Content-Type: application/json" \
    --data-binary '{"git":{...},"docker":{...},"type":"link-image"}' \
    https://webhook.atomist.com/atomist/link-image/teams/TEAM_ID
```

replacing the ellipses with the appropriate JSON and `TEAM_ID` with
the appropriate team ID.

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

```console
$ go get github.com/atomist/k8vent
```

The source code will be under `$GOPATH/src/github.com/atomist/k8vent`.
If `$GOPATH/bin` is in your `PATH`, then the `k8vent` binary will be
in your path when the above command completes successfully.  Then you
can run k8vent locally simply by invoking `k8vent` from your terminal.

If you make changes to the code, you can run tests using the Go
tooling

```console
$ go test ./...
```

or you can use `make`

```console
$ make test
```

To build, test, install, and vet, just run

```console
$ make
```

---

Created by [Atomist][atomist].
Need Help?  [Join our Slack team][slack].

[atomist]: https://atomist.com/
[slack]: https://join.atomist.com/
