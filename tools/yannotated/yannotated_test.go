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

// Bar is a struct.
type Bar struct {
	// Baz is baz.
	Baz int `yaml:"baz"`
}

func TestMain(t *testing.T) {
	tmp, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.RemoveAll(tmp.Name())

	err = mainE(Flags{
		Dir