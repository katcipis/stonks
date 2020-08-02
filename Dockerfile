ARG GOVERSION

FROM golang:${GOVERSION} as builder

WORKDIR /build
COPY . .

RUN go build -o user-manager ./cmd/user-manager/main.go

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /build/user-manager .

ENTRYPOINT ["./user-manager"]
