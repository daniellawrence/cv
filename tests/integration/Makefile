.PHONY: test

# Default target - run tests using invoke
all: test

test:
	@if [ ! -d ".venv" ]; then echo "Creating virtual environment..."; python3 -m venv .venv; fi
	source .venv/bin/activate && python -m pytest -v

help:
	@echo "Available targets:"
	@echo "  test       - Run all integration tests (default)"