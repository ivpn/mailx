# stage 1: building application binary file
FROM golang:1.24@sha256:d2d2bc1c84f7e60d7d2438a3836ae7d0c847f4888464e7ec9ba3a1339a1ee804 AS builder

RUN mkdir /daemon
ADD . /daemon
WORKDIR /daemon

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o daemon main.go

# stage 2: copy only the application binary file to a minimal, non-root container
FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d AS production

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /daemon/daemon /usr/local/bin/daemon
USER app

# run the service on container startup
EXPOSE 10025
CMD ["daemon"]
