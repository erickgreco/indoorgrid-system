MIGRATIONS_DIR=internal/db/migrations
DB_DRIVER=postgres

.PHONY: goose-up goose-down goose-status goose-create

goose-up:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" up

goose-down:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" down

goose-status:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" status

goose-create:
	@read -p "Migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql