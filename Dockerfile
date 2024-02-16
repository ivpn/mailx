# stage 1: building application binary file
FROM golang:1.22 AS builder

RUN mkdir /app
ADD . /app
WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o app cmd/main.go

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:latest AS production

COPY --from=builder /app .

# run the service on container startup
CMD ["./app"]
