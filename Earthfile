VERSION 0.8

build-configuration-image:
    FROM --platform=linux/amd64 scratch
    COPY config/docker-compose.yml /docker-compose.yml
    COPY config/gateway/Caddyfile /gateway/Caddyfile

    SAVE IMAGE --push --no-manifest-list \
    us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config

build-all:
    BUILD +build-configuration-image
    BUILD --platform=linux/amd64 ./workload+build
    BUILD --platform=linux/amd64 ./ledger+build

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
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance__short__latest

run-remote-fast:
    FROM curlimages/curl
    ARG USERNAME=formance
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance -d '{"params": {"custom.duration":"0.1", "antithesis.report.recipients":"formance-antithesis-aaaamsevowqgd236sncthxc4tu@antithesisgroup.slack.com"}}'