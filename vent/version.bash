#! /bin/bash
# calculate pre-release version name

declare Pkg=go-generate-version
declare PkgVersion=0.1.0

shopt -s extglob

function msg () {
    echo "$Pkg: $*"
}

function err () {
    msg "$*" 1>&2
}

# create version file based on VERSION environment variable
function main () {
    local version_file=v.go
    if [[ $VERSION == release ]]; then
        if ! echo -e "package vent\nconst Version = ReleaseVersion" > "$version_file"; then
            err "failed to create version file for release: $version_file"
            return 1
        fi
        return 0
    elif [[ $VERSION ]]; then
        if ! echo -e "package vent\nconst Version = \"$VERSION\"" > "$version_file"; then
            err "failed to create version file for '$VERSION': $version_file"
            return 1
        fi
        return 0
    fi

    local branch
    branch=$(git rev-parse --abbrev-ref HEAD)
    if [[ ! $branch ]]; then
        err "failed to get current branch"
        return 1
    fi
    local branch_prefix
    if [[ $branch == master ]]; then
        branch_prefix=
    else
        local safe_branch
        safe_branch=${branch//+([_\/])/-}
        branch_prefix=branch-$safe_branch.
    fi
    local ts
    ts=$(date -u '+%Y%m%d%H%M%S')
    if [[ ! $ts ]]; then
        err "failed to get current timestamp"
        return 2
    fi
    local prerelease=-${branch_prefix}$ts
    if ! echo -e "package vent\nconst Version = ReleaseVersion + \"$prerelease\"" > "$version_file"; then
        err "failed to create version file for prerelease"
        return 1
    fi
}

main "$@"
