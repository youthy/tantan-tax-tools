SHELL := /bin/bash
.PHONY: all

all: 
	go build  -o build/bin/tax_tool cmd/main.go
	@echo -e "\033[32mbuild tax_tool successfully\033[0m"

