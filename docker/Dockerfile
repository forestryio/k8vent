FROM debian:jessie

MAINTAINER David Dooling <david@atomist.com>

WORKDIR /opt/k8vent

RUN apt-get update && apt-get install -y \
        curl \
    && rm -rf /var/lib/apt/lists/*

ENV DUMB_INIT_VERSION=1.2.0

RUN curl -s -L -O https://github.com/Yelp/dumb-init/releases/download/v$DUMB_INIT_VERSION/dumb-init_${DUMB_INIT_VERSION}_amd64.deb \
    && dpkg -i dumb-init_${DUMB_INIT_VERSION}_amd64.deb \
    && rm -f dumb-init_${DUMB_INIT_VERSION}_amd64.deb

COPY ./info.json ./

COPY ./k8vent ./

ENTRYPOINT ["dumb-init", "./k8vent"]

CMD ["run"]
