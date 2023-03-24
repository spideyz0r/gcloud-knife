package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		p         string
		c         string
		expected1 bool
		expected2 error
	}{
		{p: "myproject", c: "", expected1: true, expected2: nil},
		{p: "", c: "mycommand", expected1: false, expected2: fmt.Errorf("Argument project cannot be empty")},
		{p: "", c: "", expected1: false, expected2: fmt.Errorf("Argument project cannot be empty")},
		{p: "myproject", c: "mycommand", expected1: false, expected2: nil},
	}

	for _, test := range tests {
		actual1, actual2 := parseArgs(test.p, test.c)
		if actual1 != test.expected1 {
			t.Errorf("parseArgs(%q, %q) returned %v, expected %v", test.p, test.c, actual1, test.expected1)
		}
		if (actual2 == nil) != (test.expected2 == nil) || (actual2 != nil && actual2.Error() != test.expected2.Error()) {
			t.Errorf("parseArgs(%q, %q) returned error %v, expected error %v", test.p, test.c, actual2, test.expected2)
		}
	}
}

func TestGetLastToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "abc/def/ghi", expected: "ghi"},
		{input: "foo/bar", expected: "bar"},
		{input: "bar", expected: "bar"},
		{input: "", expected: ""},
	}

	for _, test := range tests {
		actual := getLastToken(test.input)
		if actual != test.expected {
			t.Errorf("getLastToken(%q) returned %q, expected %q", test.input, actual, test.expected)
		}
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		host     string
		user     string
		expected string
	}{
		{host: "example.com", user: "", expected: os.Getenv("USER")},
		{host: "myserver", user:"", expected: "myuser"},
		{host: "myserver", user:"myuser", expected: "myuser"},
	}

	d, err := ioutil.TempDir("", "myuser")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(d)
	os.Setenv("HOME", d)

	os.MkdirAll(d + "/.ssh", 0700)
	ioutil.WriteFile(d +"/.ssh/config", []byte("Host myserver\nUser myuser\n"), 0600)

	for _, test := range tests {
		actual := getUser(test.host, test.user)
		if actual != test.expected {
			t.Errorf("getUser(%q, %q) returned %q, expected %q", test.host, test.user, actual, test.expected)
		}
	}
}

