FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o nhl .

FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/nhl /bin/nhl

EXPOSE 8080

CMD ["/bin/nhl"]