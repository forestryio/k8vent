FROM debian:buster-20200224

LABEL maintainer="Atomist <docker@atomist.com>"

RUN apt-get update && apt-get install -y \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 2866 atomist \
    && useradd --home-dir /home/atomist --create-home --uid 2866 --gid 2866 --shell /bin/sh --skel /dev/null atomist

WORKDIR /opt/k8svent

COPY ./k8svent ./

ENTRYPOINT ["./k8svent"]

USER atomist:atomist
