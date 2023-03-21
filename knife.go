package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	compute "cloud.google.com/go/compute/apiv1"
	protobuf "cloud.google.com/go/compute/apiv1/computepb"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"google.golang.org/api/iterator"

	"github.com/kevinburke/ssh_config"
)

type Knife struct {
	targets []Target
	command string
	user    string
}

type Target struct {
	zone         string
	name         string
	machine_type string
	hostname     string
}

func (k Knife) Print() {
	fmt.Printf("%s\t%s\t%s\t%s\t\n", "Hostname", "Name", "Zone", "Machine Type")
	for _, t := range k.targets {
		fmt.Printf("%s\t%s\t%s\t%s\t\n", t.hostname, t.name, t.zone, t.machine_type)
	}
}

func getUser(host string) string {
	var user string
	ssh_config_file := filepath.Join(os.Getenv("HOME"), ".ssh", "config")

	if _, err := os.Stat(ssh_config_file); errors.Is(err, os.ErrNotExist) {
		return os.Getenv("USER")
	}

	f, _ := os.Open(ssh_config_file)
	cfg, _ := ssh_config.Decode(f)
	user, err := cfg.Get(host, "User")
	if err != nil {
		log.Fatal(err)
	}
	if len(user) == 0 {
		user = os.Getenv("USER")
	}
	return user
}

func runCommand(host string, command string, user string, wg *sync.WaitGroup) {
	defer wg.Done()

	// use the private key from the ssh-agent
	sock, err := net.Dial("unix", os.Getenv(("SSH_AUTH_SOCK")))
	if err != nil {
		log.Fatal(err)
	}
	agent_client := agent.NewClient(sock)

	config := &ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(agent_client.Signers),
		},
	}

	conn, err := ssh.Dial("tcp", net.JoinHostPort(host, "22"), config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	output, err := session.Output(command)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %s\n", host, strings.Replace(string(output), "\n", " ", -1))
}

func getTarget(f string, p string) ([]Target, error) {
	var targets []Target
	ctx := context.Background()
	c, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch instances: %s\n", err)
	}
	defer c.Close()

	filter := f + " status:running"
	req := &protobuf.AggregatedListInstancesRequest{
		Project: p,
		Filter:  &filter,
	}
	it := c.AggregatedList(ctx, req)

	for {
		ret, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Failed to obtain instance: %s\n", err)
		}
		instances := ret.Value.Instances
		if len(instances) > 0 {
			for _, instance := range instances {
				targets = append(targets, Target{
					zone:         getLastToken(instance.GetZone()),
					hostname:     instance.GetHostname(),
					name:         instance.GetName(),
					machine_type: getLastToken(instance.GetMachineType()),
				})
			}
		}
	}
	return targets, nil
}

func getLastToken(s string) string {
	return strings.Split(s, "/")[len(strings.Split(s, "/"))-1]
}
