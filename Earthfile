VERSION 0.8

build:
    FROM scratch
    COPY config/docker-compose.yml /docker-compose.yml
    COPY config/volumes/database /volumes/database

    SAVE IMAGE --push us-central1-docker.pkg.dev/molten-verve-216720/formance-repository/antithesis-config:latest

push:
    BUILD +build
    BUILD ./workload+build