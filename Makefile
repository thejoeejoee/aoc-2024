include .env
export

DIRS = $(shell find . -maxdepth 1 -type d)

.PHONY: all $(DIRS)

all: $(DIRS)

$(DIRS):
	@go run $(notdir $@)/main.go