# stage 1: building application binary file
FROM golang:1.23 AS builder

RUN mkdir /api
ADD . /api
WORKDIR /api

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o api cmd/main.go

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:latest AS production

COPY --from=builder /api .

# run the service on container startup
CMD ["./api"]
