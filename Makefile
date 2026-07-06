CONFIG_LOCATION = $(HOME)/.config/danzmen/config.toml 
all: test clean watch

watch:
	air .

test:
	@echo "tested"

build:
	GOOS=linux GOARCH=amd64 go build -o danzmen .

install: 
	ln -sf $(CURDIR)/danzmen /bin/danzmen

clean:
	@echo "cleaned"

.PHONY: all watch test clean build install
