# stage 1: building application binary file
FROM golang:1.24 AS builder

RUN mkdir /daemon
ADD . /daemon
WORKDIR /daemon

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o daemon main.go

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:latest AS production

COPY --from=builder /daemon .

# run the service on container startup
EXPOSE 10025
CMD ["./daemon"]
