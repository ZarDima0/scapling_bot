include docker/.env
export $(shell sed 's/=.*//' .env)

MIGRATE=migrate
MIGRATIONS_DIR=./migrations


migrate-up:
	@echo "Applying migrations..."
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_DSN) up

migrate-create:
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	$(MIGRATE) create -ext sql -dir $(MIGRATIONS_DIR) -seq $${name}

migrate-status:
	@echo "Checking migration status..."
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_DSN) version

migrate-down:
	@echo "Rolling back migrations..."
	@$(MIGRATE) -path $(MIGRATIONS_DIR) -database $(DB_DSN) down