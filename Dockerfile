FROM alpine:latest
ENV profile=1
ADD go.mod main.go utils.go /app/
RUN apk update && apk add --no-cache go
RUN cd /app && go build && cp qb-tracker-updater /bin
CMD [ "qb-tracker-updater -conf /qBittorrent.conf -profile ${profile}" ]