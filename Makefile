# user's system
CONFIG_DIR := $(HOME)/.config/danzmen
CONFIG_LOCATION := $(CONFIG_DIR)/config.toml
LOCAL_DIRECTORY_FILES := $(HOME)/.local/share/danzmen
PREFIX ?= /usr/local

# project utils
OUTPUT_PROGRAM := $(PWD)/build/danzmen
EXAMPLE_FILE := $(PWD)/example/config.toml

config:
	mkdir -p ${CONFIG_DIR}
	test -f ${CONFIG_LOCATION} || cp ${EXAMPLE_FILE} ${CONFIG_LOCATION}
	@echo "\n >> config file location: ${CONFIG_LOCATION}\n"

watch:
	air .

test:
	@echo "tested"

build:
	GOOS=linux GOARCH=amd64 go build -o ${OUTPUT_PROGRAM} .

install: build
	install -m 755 ${OUTPUT_PROGRAM} ${PREFIX}/bin/danzmen
	mkdir -p ${CONFIG_LOCATION}
	touch ${CONFIG_LOCATION}/config.toml

symlink:
	ln -s ${CONFIG_LOCATION} $(HOME)/Desktop/danzmen.config.toml

del-db:
	rm -f ${LOCAL_DIRECTORY_FILES}/danzmen.db

clean:
	rm -f ${OUTPUT_PROGRAM} ${PREFIX}/bin/danzmen

all: test clean watch
.PHONY: all watch test clean build install del-db config symlink
