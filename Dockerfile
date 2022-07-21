FROM golang:1.18-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -o go-canal

FROM alpine

WORKDIR /app

RUN apk add --no-cache mariadb-client

COPY --from=builder /app/go-canal /app/

ENTRYPOINT ["./go-canal", "-config", "/etc/canal/config.yaml"]