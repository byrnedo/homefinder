FROM golang:1.21-alpine as builder
WORKDIR /server
ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates


COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-w -s" -o binary  ./cmd/native

# Runtime
FROM scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /server/binary /

ENTRYPOINT ["./binary"]
