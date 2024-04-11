VERSION 0.8

build:
    FROM scratch
    COPY config/docker-compose.yml /docker-compose.yml
    COPY config/volumes/database /volumes/database

    ARG TENANT_NAME=formance
    SAVE IMAGE us-central1-docker.pkg.dev/molten-verve-216720/$TENANT_NAME-repository/app