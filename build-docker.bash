#!/bin/bash
# build and push docker image

declare Pkg=build-docker
declare Version=3.4.0

set -o pipefail

declare Registry=${DOCKER_REGISTRY:-sforzando-dockerv2-local.jfrog.io}
declare BuildTarget
declare CleanTargets=
declare BuildDir=docker

function usage() {
    echo "$Pkg"
}

# exectute custom build commands before creating image
# usage: custom-build "$@"
function custom-build() {
    :
}

# informational messages to stdout
# usage: msg "message text"
function msg() {
    echo "$Pkg: $*"
}

# error messages to stderr
# usage: err "error text"
function err() {
    msg "$*" 1>&2
}

# copy files into build directory
# usage: copy SOURCE_PATH [TARGET_BASENAME]
function copy() {
    local src="$1"
    local target="$2"
    if [[ ! $src ]]; then
        err "copy: missing required argument: src"
        return 10
    fi
    local dest
    if [[ $target ]]; then
        dest="$BuildDir/$target"
    else
        local base
        base=${src##*/}
        dest="$BuildDir/$base"
    fi
    if [[ ! -e $src ]]; then
        err "copy: source file does not exist: $src"
        return 1
    fi
    local r
    [[ -d $src ]] && r="-r"
    if ! cp $r "$src" "$dest"; then
        err "copy: failed to copy $src to $dest"
        return 1
    fi
    CleanTargets="$CleanTargets $dest"
}

# remove generated/copied files from build docker directory
# usage: trap docker-clean EXIT
function docker-clean() {
    local status="$?"
    local t
    for t in $CleanTargets; do
        if ! rm -r "$t"; then
            err "failed to remove $t"
            ((status++))
        fi
    done
    exit $status
}

# return git abbreviated commit hash, append "-dirty" if uncommitted changes
# usage: commit=$(git-commit)
function git-commit() {
    local commit
    if [[ $TRAVIS == true && $TRAVIS_PULL_REQUEST == false ]]; then
        commit="$TRAVIS_COMMIT"
    else
        commit=$(git rev-parse HEAD)
        if [[ $? -ne 0 || ! $commit ]]; then
            err "failed to get current commit"
            return 1
        fi
    fi

    local abbrev_commit="${commit::7}"

    local dirty
    for diff_opt in "" "--cached"; do
        dirty=$(git diff --shortstat $diff_opt 2> /dev/null)
        if [[ $? -ne 0 ]]; then
            err "failed to determine git diff status"
            return 1
        fi
        if [[ $dirty ]]; then
            abbrev_commit="$abbrev_commit-dirty"
            break
        fi
    done

    echo "$abbrev_commit"
}

# return git branch
# usage: branch=$(git-branch)
function git-branch() {
    local branch
    if [[ $TRAVIS == true && $TRAVIS_PULL_REQUEST == false ]]; then
        branch="$TRAVIS_BRANCH"
    else
        branch=$(git rev-parse --abbrev-ref HEAD)
        if [[ $? -ne 0 || ! $branch ]]; then
            err "failed to determine git branch: $branch"
            return 1
        fi
    fi
    echo "$branch"
}

# create a version tag from git and Travis CI information
# if TRAVIS_BRANCH looks like a release tag, use it without the leading "v"
# erroring if it does not match VERSION,
# otherwise use a VERSION timestamp, abbreviated commit hash, and -dirty if there
# are uncommitted changes
# usage: tag=$(version-tag VERSION)
function version-tag() {
    local target_version="$1"
    if [[ ! $target_version ]]; then
        err "version-tag: missing required argument: VERSION"
        return 10
    fi
    shift

    if [[ $TRAVIS == true && $TRAVIS_PULL_REQUEST == false && $TRAVIS_BRANCH =~ ^v?[0-9]+\.[0-9]+\.[0-9]+(-([1-9][0-9]*|[-a-zA-Z0-9]+)(\.([1-9][0-9]*|[-a-zA-Z0-9]+))*)?$ ]]; then
        local release_tag="${TRAVIS_BRANCH#v}"
        if [[ $release_tag != $target_version && $release_tag != $target_version-* ]]; then
            err "release tag ($release_tag) does not match artifact version ($target_version)"
            return 1
        fi
        echo "$release_tag"
        return 0
    fi

    local commit
    commit=$(git-commit)
    if [[ $? -ne 0 || ! $commit ]]; then
        err "failed to get current commit"
        return 1
    fi

    local timestamp
    timestamp=$(date -u +%Y%m%d%H%M%S)
    if [[ $? -ne 0 || ! $timestamp ]]; then
        err "failed to get timestamp"
        return 1
    fi

    local version_tag="$target_version-$timestamp"
    if [[ $commit == *-dirty ]]; then
        version_tag="$version_tag-$commit"
    fi

    echo "$version_tag"
}

# generate and return the docker tag
# usage: tag=$(docker-tag VERSION)
function docker-tag() {
    local target_version="$1"
    if [[ ! $target_version ]]; then
        err "version-tag: missing required argument: VERSION"
        return 10
    fi
    shift

    local version_tag
    version_tag=$(version-tag "$target_version")
    if [[ $? -ne 0 || ! $version_tag ]]; then
        err "failed to get version tag: $version_tag"
        return 1
    fi

    local docker_tag="$Registry/$BuildTarget:$version_tag"
    echo "$docker_tag"
}

# tag the current git commit
# usage: git-tag DOCKER_TAG
function git-tag() {
    local docker_tag="$1"
    if [[ ! $docker_tag ]]; then
        err "git-tag: missing required argument: DOCKER_TAG"
        return 10
    fi
    shift

    if [[ $TRAVIS != true || $TRAVIS_PULL_REQUEST != false || $TRAVIS_BRANCH != master || $docker_tag == *-dirty ]]; then
        msg "not tagging commit"
        return 0
    fi

    if ! git config --global user.email "travis-ci@atomist.com"; then
        err "failed to set git config email"
        return 1
    fi
    if ! git config --global user.name "Travis CI"; then
        err "failed to set git config name"
        return 1
    fi

    local tag=${docker_tag##*:}
    local build_tag="$tag+travis$TRAVIS_BUILD_NUMBER"
    local t tags
    for t in "$tag" "$build_tag"; do
        if [[ $t != $TRAVIS_BRANCH ]]; then
            if ! git tag "$t" -a -m "Generated tag from TravisCI build $TRAVIS_BUILD_NUMBER"; then
                err "failed to tag commit: $t"
                return 1
            fi
            tags="$tags $t"
        fi
    done

    local origin=origin
    if [[ $GITHUB_TOKEN ]]; then
        origin=https://$GITHUB_TOKEN:x-oauth-basic@github.com/$TRAVIS_REPO_SLUG.git
    fi
    if ! git push "$origin" --tags; then
        err "failed to push tags: $tags"
        return 1
    fi

    msg "pushed git tags: $tags"
}

# create file with project and commit info
# usage: info-json
function info-json() {
    local target_version="$1"
    if [[ ! $target_version ]]; then
        err "info-json: missing required argument: VERSION"
        return 10
    fi
    shift

    local app_version
    app_version=$(version-tag "$target_version")
    if [[ $? -ne 0 || ! $app_version ]]; then
        err "failed to get version: $app_version"
        return 1
    fi

    local git_commit_id
    git_commit_id=$(git-commit)
    if [[ $? -ne 0 || ! $git_commit_id ]]; then
        err "failed to get commit hash: $git_commit_id"
        return 1
    fi

    local git_branch
    git_branch=$(git-branch)
    if [[ $? -ne 0 || ! $git_branch ]]; then
        err "failed to get branch: $git_branch"
        return 1
    fi

    local git_commit_time
    if [[ $git_commit_id == *-dirty ]]; then
        git_commit_time=$(date -u +%Y-%m-%dT%H:%M:%S%z)
        if [[ $? -ne 0 || ! $git_commit_time ]]; then
            err "failed to determine current time: $git_commit_time"
            return 1
        fi
    else
        local commit_utime
        commit_utime=$(git show -s --format=format:%ct "$git_commit_id")
        if [[ $? -ne 0 || ! $commit_utime ]]; then
            err "failed to determine git commit $git_commit_id time: $commit_utime"
            return 1
        fi
        local sys=$(uname -s)
        local date_exit
        case "$sys" in
            Darwin)
                git_commit_time=$(date -j -u -f %s "$commit_utime" +%Y-%m-%dT%H:%M:%S%z)
                date_exit="$?"
                ;;
            Linux)
                git_commit_time=$(date -u --date="@$commit_utime" -Iseconds)
                date_exit="$?"
                ;;
            *)
                err "unsupported system: $sys"
                return 1
                ;;
        esac
        if [[ $date_exit -ne 0 || ! $git_commit_time ]]; then
            err "failed to determine git commit time: $git_commit_time"
            return 1
        fi
    fi

    local info="$BuildDir/info.json"
    cat > "$info" <<EOF
{
  "app": {
    "version": "$app_version",
    "group_id": "github.com/atomisthq",
    "artifact_id": "$BuildTarget"
  },
  "git": {
    "branch": "$git_branch",
    "commit": {
      "id": "$git_commit_id",
      "time": "$git_commit_time"
    }
  }
}
EOF
    CleanTargets="$CleanTargets $info"
}

# perform tasks before docker build
function pre-build() {
    local build_type="$1"
    if [[ ! $build_type ]]; then
        err "pre-build: missing required argument BUILD_TYPE"
        return 10
    fi
    shift
    local target_version="$1"
    if [[ ! $target_version ]]; then
        err "pre-build: missing required argument TARGET_VERSION"
        return 10
    fi
    shift

    case "$build_type" in
        lein | mvn)
            copy "target/$BuildTarget-$target_version.jar" "$BuildTarget.jar" || return 1
            copy target/lib || return 1
            ;;
        go)
            local kernel machine
            kernel=$(uname -s)
            machine=$(uname -m)
            if [[ $kernel != Linux || $machine != x86_64 ]]; then
                msg "building $BuildTarget for docker"
                if ! GOOS=linux GOARCH=amd64 go build -v; then
                    err "failed to build $BuildTarget for linux amd64"
                    return 1
                fi
                CleanTargets="$CleanTargets $BuildTarget"
            fi
            copy "$BuildTarget" || return 1
            ;;
        js)
            copy package.json || return 1
            ;;
        sh)
            if ! install "$BuildTarget".*sh "$BuildTarget"; then
                err "failed to create executable version of $BuildTarget"
                return 1
            fi
            CleanTargets="$CleanTargets $BuildTarget"
            copy "$BuildTarget" || return 1
            ;;
        exe)
            copy "$BuildTarget" || return 1
            ;;
    esac
}

# build the docker image, CUSTOM_ARGS are passed to custom-build
# usage: docker-build CUSTOM_ARGS
function docker-build() {
    local build_type target_version
    if [[ -f project.clj ]]; then
        build_type=lein
        target_version=$(awk '$1 ~ /defproject/ { v=$3; gsub(/"/, "", v); print v; exit 0 } END { if (v == "") exit 1 }' project.clj)
    elif [[ -f pom.xml ]]; then
        build_type=mvn
        target_version=$(awk '/<version>/ { v=$1; gsub(/<\/?version>/, "", v); print v; exit 0 } END { if (v == "") exit 1 }' pom.xml)
    elif ls *.go > /dev/null 2>&1; then
        build_type=go
        target_version=$(./"$BuildTarget" version | sed "s/$BuildTarget  *//")
    elif [[ -f package.json ]]; then
        build_type=js
        target_version=$(jq -er .version package.json)
    elif ls "$BuildTarget".*sh > /dev/null 2>&1; then
        build_type=sh
        target_version=$(bash "$BuildTarget".*sh --version | sed "s/$BuildTarget  *//")
    elif [[ -x $BuildTarget ]]; then
        build_type=exe
        target_version=$(./"$BuildTarget" --version | sed "s/$BuildTarget  *//")
    elif [[ -f VERSION ]]; then
        build_type=none
        target_version=$(< VERSION)
    else
        build_type=none
        target_version=$(echo "0.0.0")
    fi
    if [[ $? -ne 0 || ! $target_version ]]; then
        err "failed to determine $BuildTarget version: $target_version"
        return 1
    fi

    custom-build "$@" || return 1

    pre-build "$build_type" "$target_version" || return 1

    info-json "$target_version" || return 1

    local tag
    tag=$(docker-tag "$target_version")
    if [[ $? -ne 0 || ! $tag ]]; then
        err "failed to create docker tag"
        return 1
    fi

    if ! (cd "$BuildDir" && docker build -t "$tag" .); then
        err "failed to build docker image"
        return 1
    fi

    if [[ $tag == *-dirty ]]; then
        msg "not pushing image or tagging commit from dirty repo: $tag"
        return 0
    fi

    msg "checking if image already exists in registry"
    if docker pull "$tag" > /dev/null 2>&1; then
        err "image $tag already exists"
        return 1
    fi

    msg "pushing $tag"
    if ! docker push "$tag"; then
        err "failed to push docker image"
        return 1
    fi
    msg "pushed artifact $tag"

    if ! git-tag "$tag"; then
        err "failed to tag commit"
        return 1
    fi
}

function main() {
    BuildTarget=${PWD##*/}
    BuildTarget=${BuildTarget#docker-}
    if [[ ! $BuildTarget ]]; then
        err "failed to determine current directory name: $PWD"
        return 1
    fi

    trap docker-clean EXIT

    docker-build "$@" || return 1
}

main "$@" || exit 1
exit 0
