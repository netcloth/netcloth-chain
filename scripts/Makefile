all: get_tools

########################################

GOLINT = github.com/tendermint/lint/golint
STATIK = github.com/rakyll/statik

GOLINT_CHECK := $(shell command -v golint 2> /dev/null)
STATIK_CHECK := $(shell command -v statik 2> /dev/null)



check_tools:
ifndef STATIK_CHECK
	@echo "No statik in path.  Install with 'make get_tools'."
else
	@echo "Found statik in path."
endif


check_dev_tools:
	$(MAKE) check_tools
ifndef GOLINT_CHECK
	@echo "No golint in path.  Install with 'make get_dev_tools'."
else
	@echo "Found golint in path."
endif

get_tools:
ifdef STATIK_CHECK
	@echo "Statik is already installed.  Run 'make update_tools' to update."
else
	@echo "Installing statik"
	go version
	go get -v $(STATIK)
endif

get_dev_tools:
	$(MAKE) get_tools
ifdef GOLINT_CHECK
	@echo "Golint is already installed.  Run 'make update_tools' to update."
else
	@echo "Installing golint"
	go get -v $(GOLINT)
endif