ARG GOVERSION

FROM golang:${GOVERSION}

WORKDIR /app

COPY . .

RUN go build -o users-manager ./cmd/users-manager/main.go

ENTRYPOINT ["/app/users-manager"]
