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
	// So useful.
	Bar Bar `yaml:"bar"`
	// Rebar is another bar.
	Rebar Bar `yaml:"rebar"`
	Quz   map[string]string
}

// Bar is a str