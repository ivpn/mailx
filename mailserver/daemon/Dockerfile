# stage 1: building application binary file
FROM golang:1.27.1@sha256:1e93e00a31255c07e9a34c4207f3006e1501730c5323697cee7dfb827fdae44c AS builder

RUN mkdir /daemon
ADD . /daemon
WORKDIR /daemon

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o daemon main.go

# stage 2: copy only the application binary file to a minimal, non-root container
FROM alpine:3.24.2@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6 AS production

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /daemon/daemon /usr/local/bin/daemon
USER app

# run the service on container startup
EXPOSE 10025
CMD ["daemon"]
