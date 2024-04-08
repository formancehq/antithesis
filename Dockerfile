FROM scratch
COPY config/docker-compose.yml /docker-compose.yml
ADD config/volumes/database /volumes/database
