package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pborman/getopt"
)

func main() {

	help := getopt.BoolLong("help", 'h', "display this help")
	filter := getopt.StringLong("filter", 'f', "", "target filter")
	command := getopt.StringLong("command", 'c', "", "command to be run externally")
	project := getopt.StringLong("project", 'p', "", "GCP project")
	user := getopt.StringLong("user", 'u', "", "[user]")

	getopt.Parse()

	if *help {
		getopt.Usage()
		os.Exit(0)
	}

	search, err := parseArgs(*project, *command)
	if err != nil {
		getopt.Usage()
		log.Fatal(err)
	}

	var k Knife
	t, err := getTarget(*filter, *project)
	if err != nil {
		log.Fatal(err)
	}

	k = Knife{
		targets: t,
		user:    *user,
		command: *command,
	}

	if search {
		k.Print()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	for _, t := range k.targets {
		wg.Add(1)
		if len(*user) == 0 {
			k.user = getUser(t.hostname)
		}
		go runCommand(t.hostname, k.command, k.user, &wg)
	}
	wg.Wait()
	os.Exit(0)
}

func parseArgs(p, c string) (bool, error) {
	if len(p) == 0 {
		return false, fmt.Errorf("Argument project cannot be empty")
	}
	if len(c) == 0 {
		return true, nil
	}
	return false, nil
}
