#!/bin/bash
# build and test a go package

set -o pipefail

declare Pkg=travis-build-go
declare Version=0.1.0

# write message to standard out (stdout)
# usage: msg MESSAGE
function msg() {
    echo "$Pkg: $*"
}

# write message to standard error (stderr)
# usage: err MESSAGE
function err() {
    msg "$*" 1>&2
}

# git tag and push
# usage: git-tag TAG[...]
function git-tag () {
    if [[ ! $1 ]]; then
        err "git-tag: missing required argument: TAG"
        return 10
    fi

    if ! git config --global user.email "travis-ci@atomist.com"; then
        err "failed to set git user email"
        return 1
    fi
    if ! git config --global user.name "Travis CI"; then
        err "failed to set git user name"
        return 1
    fi
    local tag
    for tag in "$@"; do
        if ! git tag "$tag" -m "Generated tag from Travis CI build $TRAVIS_BUILD_NUMBER"; then
            err "failed to create git tag: '$tag'"
            return 1
        fi
    done
    local remote=origin
    if [[ $GITHUB_TOKEN ]]; then
        remote=https://$GITHUB_TOKEN:x-oauth-basic@github.com/$TRAVIS_REPO_SLUG.git
    fi
    if ! git push --quiet "$remote" "$@" > /dev/null 2>&1; then
        err "failed to push git tag(s): $*"
        return 1
    fi
}

# create and echo a prerelease timestamped, and optionally branched, version
# usage: prerelease_version=$(prerelease-version BASE_VERSION [BRANCH])
function prerelease-version () {
    local base_version=$1
    if [[ ! $base_version ]]; then
        err "prerelease-version: missing required argument: BASE_VERSION"
        return 10
    fi
    shift
    local branch=$1 prerelease=
    if [[ $branch && $branch != master ]]; then
        shift
        local safe_branch
        safe_branch=$(echo -n "$branch" | tr -C -s '[:alnum:]-' . | sed -e 's/^[-.]*//' -e 's/[-.]*$//')
        if [[ $? -ne 0 || ! $safe_branch ]]; then
            err "failed to create safe branch name from '$branch': '$safe_branch'"
            return 1
        fi
        prerelease=$safe_branch.
    fi

    local timestamp
    timestamp=$(date -u +%Y%m%d%H%M%S)
    if [[ $? -ne 0 || ! $timestamp ]]; then
        err "failed to generate timestamp"
        return 1
    fi

    echo "$base_version-$prerelease$timestamp"
}

# usage: main "$@"
function main () {

    local target=${TRAVIS_REPO_SLUG##*/}
    if [[ ! $target ]]; then
        err "failed to determine targer from repo slug: '$TRAVIS_REPO_SLUG'"
        return 1
    fi

    msg "running make"
    if ! make TARGET="$target"; then
        err "make failed"
        return 1
    fi

    [[ $TRAVIS_PULL_REQUEST == false ]] || return 0

    local tag_version git_tag=
    if [[ $TRAVIS_TAG =~ ^[0-9]+\.[0-9]+\.[0-9]+(-(m|rc)\.[0-9]+)?$ ]]; then
        tag_version=$TRAVIS_TAG
    else
        local target_version
        target_version=$("./$target" version | awk 'NR == 1 { print $2 }')
        if [[ $? -ne 0 || ! $target_version ]]; then
            err "failed to determine version running './$target version': '$target_version'"
            return 1
        fi
        local prerelease_version
        prerelease_version=$(prerelease-version "$target_version" "$TRAVIS_BRANCH")
        if [[ $? -ne 0 || ! $prerelease_version ]]; then
            err "failed to create prerelease version: '$prerelease_version'"
            return 1
        fi
        tag_version=$prerelease_version
        git_tag=$tag_version
    fi

    if [[ $DOCKER_USER ]]; then
        if ! docker login -u "$DOCKER_USER" -p "$DOCKER_PASSWORD" $DOCKER_REGISTRY; then
            err "failed to login to docker registry '$DOCKER_REGISTRY'"
            return 1
        fi
        local docker_tag=
        if [[ $DOCKER_REGISTRY ]]; then
            docker_tag=$DOCKER_REGISTRY/
        fi
        docker_tag=$docker_tag$TRAVIS_REPO_SLUG:$tag_version
        msg "running make docker"
        if ! make docker DOCKER_TAG="$docker_tag"; then
            err "make docker failed: DOCKER_TAG='$docker_tag'"
            return 1
        fi
    fi

    if ! git-tag $git_tag "$tag_version+travis.$TRAVIS_BUILD_NUMBER"; then
        err "failed to tag commit"
        return 1
    fi
}

main "$@" || exit 1
exit 0
