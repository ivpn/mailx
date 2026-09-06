# stage 1: building application binary file
FROM golang:1.25@sha256:699337d620559a59b4a2bb298ad59611e535d2ee755a34cf2d2a98f37578dc80 AS builder

RUN mkdir /api
ADD . /api
WORKDIR /api

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o api ./cmd/

# stage 2: copy only the application binary file to a minimal, non-root container
FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d AS production

RUN addgroup -S app && adduser -S app -G app
# owned so the /var/log named volume inherits this ownership on first mount
RUN mkdir -p /var/log && chown app:app /var/log
COPY --from=builder /api/api /usr/local/bin/api
USER app

# run the service on container startup
CMD ["api"]
