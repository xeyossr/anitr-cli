GO=go
BINARY_NAME=anitr-cli
INSTALL_DIR=/usr/bin
VERSION=v4.1.0

.PHONY: mod-tidy build run install clean all test dev-build release help

mod-tidy:
	$(GO) mod tidy

build: mod-tidy
	$(GO) build -o $(BINARY_NAME)

dev-build: mod-tidy
	$(GO) build -race -o $(BINARY_NAME)-dev

release: mod-tidy
	$(GO) build -ldflags "-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

dev-run: dev-build
	./$(BINARY_NAME)-dev

test:
	$(GO) test ./...

test-verbose:
	$(GO) test -v ./...

install: build
	chmod +x $(BINARY_NAME)
	sudo mv $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)

uninstall:
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME)-dev

config-backup:
	@if [ -d ~/.config/anitr-cli ]; then \
		cp -r ~/.config/anitr-cli ~/.config/anitr-cli-backup-$(shell date +%Y%m%d-%H%M%S); \
		echo "Config backed up to ~/.config/anitr-cli-backup-$(shell date +%Y%m%d-%H%M%S)"; \
	else \
		echo "No config directory found"; \
	fi

config-restore:
	@echo "Available backups:"
	@ls -la ~/.config/anitr-cli-backup-* 2>/dev/null || echo "No backups found"

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  dev-build    - Build with race detection for development"
	@echo "  release      - Build optimized release version"
	@echo "  run          - Build and run the application"
	@echo "  dev-run      - Build and run development version"
	@echo "  test         - Run tests"
	@echo "  test-verbose - Run tests with verbose output"
	@echo "  install      - Install to system"
	@echo "  uninstall    - Remove from system"
	@echo "  clean        - Clean build artifacts"
	@echo "  config-backup- Backup user configuration"
	@echo "  config-restore- Show available config backups"
	@echo "  mod-tidy     - Tidy Go modules"
	@echo "  help         - Show this help"

all: mod-tidy build install
