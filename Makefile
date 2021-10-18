
.PHONY: default build docker-image test stop clean-images clean

BINARY = kubewatch

VERSION=
BUILD=

PKG            = github.com/bitnami-labs/kubewatch
TRAVIS_COMMIT ?= `git describe --tags`
GOCMD          = go