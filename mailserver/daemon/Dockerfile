# stage 1: building application binary file
FROM golang:1.27.3@sha256:e519d3b0a5b285e5aa53ab9eef02bc6bbb9cb04608389b743a136a9ce75f0d1e AS builder

RUN mkdir /daemon
ADD . /daemon
WORKDIR /daemon

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o daemon main.go

# stage 2: copy only the application binary file to a minimal, non-root container
FROM alpine:3.24.3@sha256:be2656061b536c1883c4034790e21d0a09c5aaaf51968b89368a3c0371ea2c9a AS production

RUN addgroup -S app && adduser -S app -G app
COPY --from=builder /daemon/daemon /usr/local/bin/daemon
USER app

# run the service on container startup
EXPOSE 10025
CMD ["daemon"]
