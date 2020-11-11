FROM golangci/golangci-lint:v1.32.2 as build

WORKDIR /build

COPY ./ ./

RUN make

FROM debian:buster-20201012

LABEL maintainer="Atomist <docker@atomist.com>"

RUN apt-get update && apt-get install -y \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 2866 atomist \
    && useradd --home-dir /home/atomist --create-home --uid 2866 --gid 2866 --shell /bin/sh --skel /dev/null atomist

WORKDIR /opt/k8svent

ENTRYPOINT ["./k8svent"]

COPY --from=build /build/k8svent ./

USER atomist:atomist
