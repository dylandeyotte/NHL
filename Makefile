
migrate:
	go run github.com/pressly/goose/v3/cmd/goose@latest \
	-dir ./sql/schema \
	postgres "$(DB_URL)" up

