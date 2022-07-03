NAME=$(shell basename $(CURDIR))
VER?="0.0.0-dev"
VERSION=$(shell echo $(VER) | sed -e 's/^v//g' -e 's/\//_/g')

ifeq ($(shell uname -m),x86_64)
  ARCH=x86_64
else ifeq ($(shell uname -m),arm64)
  ARCH=aarch64
endif


.PHONY: help
help: ## Show this help.
	@IFS=$$'\n' ; \
	lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//'`); \
	for line in $${lines[@]}; do \
		IFS=$$'#' ; \
		split=($$line) ; \
		command=`echo $${split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		info=`echo $${split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-38s %s\n" $$command $$info ; \
	done


.PHONY: printvars
printvars:  ## print all variables
	@$(foreach V,$(sort $(.VARIABLES)), \
	   $(if $(filter-out environment% default automatic, \
	   $(origin $V)),$(info $V=$($V) ($(value $V)))))


.PHONY: get-version
get-version:  ## Return version
	@echo $(VERSION) > VERSION
	@echo $(VERSION)


.PHONY: tests
tests:  dev-env-up check  ## Spin environment and run nix flake check


.PHONY: check
check:  ## Run nix flake check
	nix flake check --print-build-logs


.PHONY: build
build:  ## Build application and places the binary under ./result/bin
	@echo $(VERSION) > VERSION
	./build/common/nix.sh build --print-build-logs


.PHONY: build-docker-image
build-docker-image:  ## Build docker container for native architecture
	@echo $(VERSION) > VERSION
	./build/common/nix-docker-image.sh
	docker load < result


.PHONY: dev-env-up
dev-env-up: dev-env-down _dev-env-build _dev-env-up ## Starts development environment


.PHONY: dev-env-down
dev-env-down: _dev-env-down  ## Stops development environment


.PHONY: dev-env-build
dev-env-build: _dev-env-build  ## Builds development environment
