VERSION 0.8

build-configuration-image:
    FROM scratch
    COPY config/docker-compose.yml /docker-compose.yml
    RUN mkdir /volumes/database

    SAVE IMAGE --push us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config:latest

build-all:
    BUILD +build-configuration-image
    BUILD ./workload+build

run:
    LOCALLY
    RUN earthly ./workload+build
    RUN --no-cache rm -rf config/volumes/database/*
    RUN --no-cache docker compose -f config/docker-compose.yml up workload # Wait to let the database starting property (todo: need to add this on the ledger maybe)
    RUN --no-cache docker compose -f config/docker-compose.yml down -v