SHELL := /bin/bash
.PHONY: debug testbed clean init
debug:
	source env.sh && enter_goland
testbed:
	source env.sh && clean && start
clean:
	source env.sh && clean
init:
	bash entrypoint.sh

