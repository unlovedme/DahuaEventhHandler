package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// Config is a config.
type Config struct {
	// Foo is foo.
	Foo string `yaml:"foo"`
	// Bar is bar.
	// So