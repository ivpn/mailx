# stage 1: building application binary file
FROM golang:1.27.4@sha256:205c7636d74e0ffbe3d007162df7d74f859f5a8c8322dfb9f907977686c126cf AS builder

RUN mkdir /api
ADD . /api
WORKDIR /api

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

RUN go build -o api ./cmd/

# stage 2: copy only the application binary file to a minimal, non-root container
FROM alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf AS production

RUN addgroup -S app && adduser -S app -G app
# owned so the /var/log named volume inherits this ownership on first mount
RUN mkdir -p /var/log && chown app:app /var/log
COPY --from=builder /api/api /usr/local/bin/api
USER app

# run the service on container startup
CMD ["api"]
