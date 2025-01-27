.PHONY: db import create delete

help:
	@echo "Usage: make [subcommand]"
	@echo "Available subcommands:"
	@echo "  db create       - Create resources in a container"
	@echo "  db delete       - Delete resources from a container"
	@echo "  import          - Import CSV into a container"
	@echo "  help            - Show this help message"

db:
	@echo "Usage: make db [subcommand] CONTAINER_ID=<container_id>"
	@echo "Available subcommands for db:"
	@echo "  create    - Create resources in the container"
	@echo "  delete    - Delete resources from the container"

create:
	@if [ -z "$(CONTAINER_ID)" ]; then \
		echo "Error: CONTAINER_ID is required"; \
		exit 1; \
	fi
	@echo "Creating resources in container $(CONTAINER_ID)..."
	docker exec -it $(CONTAINER_ID) go run main.go db create

delete:
	@if [ -z "$(CONTAINER_ID)" ]; then \
		echo "Error: CONTAINER_ID is required"; \
		exit 1; \
	fi
	@echo "Deleting resources from container $(CONTAINER_ID)..."
	docker exec -it $(CONTAINER_ID) go run main.go db delete

import:
	@if [ -z "$(CSV_PATH)" ] || [ -z "$(CONTAINER_ID)" ]; then \
		echo "Error: CSV_PATH and CONTAINER_ID are required"; \
		exit 1; \
	fi
	@echo "Copying CSV to container $(CONTAINER_ID)..."
	docker cp $(CSV_PATH) $(CONTAINER_ID):/tmp/$(notdir $(CSV_PATH))
	@echo "Importing CSV into container $(CONTAINER_ID)..."
	docker exec -it $(CONTAINER_ID) go run main.go import /tmp/$(notdir $(CSV_PATH))
	@echo "Cleaning up CSV in container..."
	docker exec -it $(CONTAINER_ID) rm -f /tmp/$(notdir $(CSV_PATH))
	@echo "CSV import completed."

