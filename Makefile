## the publicly-visible name of your binary
NAME=go-vertica

## the go-get'able path
PKG_PATH=github.com/tobz/$(NAME)

## version, taken from Git tag (like v1.0.0) or hash
VER:=$(shell (git describe --always --dirty 2>/dev/null || echo "¯\\\\\_\\(ツ\\)_/¯") | sed -e 's/^v//g' )

BIN=.godeps/bin
GPM=$(BIN)/gpm
GPM_LINK=$(BIN)/gpm-link
GVP=$(BIN)/gvp

## @todo should use "$(GVP) in", but that fails
## all non-test source files
SOURCES:=$(shell go list -f '{{range .GoFiles}}{{ $$.Dir }}/{{.}} {{end}}' ./... | sed -e "s@$(PWD)/@@g" )

## all packages in this prject
PACKAGES:=$(shell go list -f '{{.Name}}' ./... )

.PHONY: all devtools deps test build clean rpm

## targets after a | are order-only; the presence of the target is sufficient
## http://stackoverflow.com/questions/4248300/in-a-makefile-is-a-directory-name-a-phony-target-or-real-target

all: build

$(BIN) stage:
	mkdir -p $@

$(GPM): | $(BIN)
	curl -s -L -o $@ https://github.com/pote/gpm/raw/v1.3.1/bin/gpm
	chmod +x $@

$(GPM_LINK): | $(BIN)
	curl -s -L -o $@ https://github.com/elcuervo/gpm-link/raw/v0.0.1/bin/gpm-link
	chmod +x $@

$(GVP): | $(BIN)
	curl -s -L -o $@ https://github.com/pote/gvp/raw/v0.1.0/bin/gvp
	chmod +x $@

.godeps/.gpm_installed: $(GPM) $(GVP) $(GPM_LINK) Godeps
	test -e .godeps/src/$(PKG_PATH) || $(GVP) in $(GPM) link add $(PKG_PATH) $(PWD)
	$(GVP) in $(GPM) install
	touch $@

$(BIN)/ginkgo: .godeps/.gpm_installed
	$(GVP) in go install github.com/onsi/ginkgo/ginkgo
	touch $@

$(BIN)/mockery: .godeps/.gpm_installed
	$(GVP) in go install github.com/vektra/mockery
	touch $@

## installs dev tools
devtools: $(BIN)/ginkgo $(BIN)/mockery

## just installs dependencies
deps: .godeps/.gpm_installed

## run tests
test: $(BIN)/ginkgo
	$(GVP) in $(BIN)/ginkgo $(PACKAGES)

## build the binary
## augh!  gvp shell escaping!!
## https://github.com/pote/gvp/issues/22
stage/$(NAME): .godeps/.gpm_installed $(SOURCES) | stage
	$(GVP) in go build -o $@ -ldflags '-X\ main.version\ $(VER)' -v .

stage/vertica_test: .godeps/.gpm_installed $(SOURCES) | stage
	$(GVP) in go build -o stage/vertica_test -v vertica_test/main.go

## same, but shorter
build: test stage/$(NAME)

## duh
clean:
	rm -rf stage .godeps release

rpm: build
	mkdir -p stage/rpm/usr/bin stage/rpm/etc
	
	cp stage/$(NAME) stage/rpm/usr/bin/
	chmod 555 stage/rpm/usr/bin/$(NAME)
	
	## config files
	cp etc/* stage/rpm/etc/

	cd stage && fpm \
	    -s dir \
	    -t rpm \
	    -n $(NAME) \
	    -v $(VER) \
	    --rpm-use-file-permissions \
	    -C rpm \
	    etc usr
