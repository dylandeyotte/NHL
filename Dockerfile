FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN CGO_ENABLED=0 go build -o nhl .
RUN CGO_ENABLED=0 go build -o load_teams ./cmd/load_teams

FROM debian:stable-slim

RUN apt-get update && \
    apt-get install -y ca-certificates tzdata && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/nhl /bin/nhl
COPY --from=builder /app/sql/schema /sql/schema
COPY --from=builder /app/load_teams /bin/load_teams
COPY --from=builder /go/bin/goose /bin/goose

EXPOSE 8080

CMD /bin/goose -dir /sql/schema postgres "$DB_URL" up && \
/bin/load_teams && \
/bin/nhl