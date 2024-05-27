# builder image
FROM golang AS build

WORKDIR /app

ADD go.mod main.go /app/
RUN go build -o qb-tracker-updater main.go

# base image
FROM cgr.dev/chainguard/wolfi-base
RUN apk add ca-certificates

COPY --from=build /app/qb-tracker-updater /bin

VOLUME /config

ENV PROFILE=1
CMD [ "/bin/qb-tracker-updater", "--conf", "/config/qBittorrent.conf", "--profile ${PROFILE}" ]