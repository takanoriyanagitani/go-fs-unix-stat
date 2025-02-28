#!/bin/sh

ls \
	./main.go \
	./run.sh |
	./stat2json |
	jq -c
