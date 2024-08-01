VERSION 0.8

build-configuration-image-001:
    FROM --platform=linux/amd64 scratch
    COPY config-001/docker-compose.yml /docker-compose.yml
    COPY config-001/gateway/Caddyfile /gateway/Caddyfile

    SAVE IMAGE --push --no-manifest-list \
    us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config-001:latest

build-configuration-image-002:
    FROM --platform=linux/amd64 scratch
    COPY config-002/docker-compose.yml /docker-compose.yml
    COPY config-002/gateway/Caddyfile /gateway/Caddyfile

    SAVE IMAGE --push --no-manifest-list \
    us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config-002:latest

build-all:
    BUILD +build-configuration-image-001
    BUILD +build-configuration-image-002
    BUILD --platform=linux/amd64 ./workload+build
    BUILD --platform=linux/amd64 ./ledger+build

run:
    LOCALLY
    RUN earthly ./workload+build
    RUN earthly ./ledger+build
    RUN --no-cache rm -rf config-001/volumes/database/*
    RUN --no-cache docker compose -f config-001/docker-compose.yml up workload
    RUN --no-cache docker compose -f config-001/docker-compose.yml down -v

run-remote:
    FROM curlimages/curl
    ARG USERNAME=formance
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance__short__latest

run-remote-fast:
    FROM curlimages/curl
    ARG USERNAME=formance
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance -d '{"params": {"custom.duration":"0.1", "antithesis.report.recipients":"formance-antithesis-aaaamsevowqgd236sncthxc4tu@antithesisgroup.slack.com"}}'