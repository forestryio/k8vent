# k8vent

[![Build Status](https://travis-ci.org/atomist/k8vent.svg?branch=master)](https://travis-ci.org/atomist/k8vent)

Send kubernetes pods as JSON to webhook endpoints.  k8vent is
typically run from its Docker image in a kubernetes cluster to send
pod state changes to the Atomist kubernetes webhook endpoint for your
team.

## Running

You can use [k8vent-deployment.json][k8vent-deployment] as a starting
point for running k8vent in your kubernetes cluster.  Be sure to
update the value of the `K8VENT_WEBHOOKS` environment variable,
replacing the last element of the URL with the ID of your Slack team.
You can get your team ID from https://app.atomist.com/teams or by
sending `team` as a message to the Atomist bot, e.g., `@atomist team`,
in Slack.

You can optionally change the value of the
`ATOMIST_ENVIRONMENT` environment variable to a meaningful name for
your kubernetes cluster, e.g., "production" or "staging".

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

```
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
