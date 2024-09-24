VERSION 0.8

build-all:
    BUILD ./config-001/+build
    BUILD ./config-002/+build
    BUILD --platform=linux/amd64 ./workload+build
    BUILD --platform=linux/amd64 ./ledger+build

run-local:
    LOCALLY
    RUN --no-cache rm -rf config-001/volumes/database/*
    RUN --no-cache docker compose -f config-001/docker-compose.yml up workload
    RUN --no-cache docker compose -f config-001/docker-compose.yml down -v

run-remote:
    ARG USERNAME=formance
    ARG EXPERIMENT
    FROM curlimages/curl
    COPY experiments /experiments
    RUN --no-cache --secret ANTITHESIS_PASSWORD curl --fail --user "$USERNAME:$ANTITHESIS_PASSWORD" -X POST https://formance.antithesis.com/api/v1/launch_experiment/formance -d @/experiments/$EXPERIMENT.json

start-debugging:
    ARG SESSION_ID
    ARG INPUT_HASH
    ARG VTIME
    ARG EMAIL
    FROM alpine/httpie
    RUN --no-cache --secret ANTITHESIS_PASSWORD http --ignore-stdin -a "formance:$ANTITHESIS_PASSWORD" post https://formance.antithesis.com/api/v1/launch_experiment/launch_debugging \
    "params.antithesis.debugging.session_id"=$SESSION_ID \
    "params.antithesis.debugging.input_hash"=$INPUT_HASH \
    "params.antithesis.debugging.vtime"=$VTIME \
    "params.antithesis.report.recipients"=$EMAIL