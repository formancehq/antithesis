VERSION 0.8

build-configuration-image:
    FROM --platform=linux/amd64 scratch
    COPY config/docker-compose.yml /docker-compose.yml
    COPY config/gateway/Caddyfile /gateway/Caddyfile

    SAVE IMAGE --push us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config:latest

build-all:
    BUILD +build-configuration-image
    BUILD ./workload+build
    BUILD ./ledger+build

run:
    LOCALLY
    RUN earthly ./workload+build
    RUN earthly ./ledger+build
    RUN --no-cache rm -rf config/volumes/database/*
    RUN --no-cache docker compose -f config/docker-compose.yml up workload
    RUN --no-cache docker compose -f config/docker-compose.yml down -v

run-remote:
    FROM curlimages/curl
    ARG USERNAME=formance
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance__configuration__latest