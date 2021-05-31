PLUGIN_NAME="HashAnnotator"

default: install

.PHONY: install
install: go-install clean build install-plugin

.PHONY: install-plugin
install-plugin:
	./install.sh

.PHONY: build
build:
	go build -o $(PLUGIN_NAME)

.PHONY: clean
clean:
	rm $(PLUGIN_NAME) || true
	rm -rf $(XDG_CONFIG_HOME)/kustomize/plugin/pcjun97/v1/hashannotator || true

.PHONY: go-install
go-install:
	go install
