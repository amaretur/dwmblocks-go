
BIN=dwmblocks-go
CONFIG=config.json

BIN_DIR=/usr/bin
CONFIG_DIR=$(HOME)/.config/dwmblocks

all: install

build:
	go build -o ${BIN} .

install: build
	sudo mv ${BIN} ${BIN_DIR}
	mkdir -p ${CONFIG_DIR}
	cp ${CONFIG} ${CONFIG_DIR}

uninstall:
	sudo rm -f ${BIN_DIR}/${BIN}
	rm -f ${CONFIG_DIR}/${CONFIG}

.PHONY: install, uninstall
