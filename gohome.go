package main

import (
	"github.com/urfave/cli"
	"os"
)

var conf *Conf

func main() {
	var server *Server = new(Server)
	var client *Client = new(Client)
	conf = &Conf{
		Server: server,
		Client: client,
	}

	app := cli.NewApp()
	app.Name = "GoHome"
	app.Description = "GoHome is a dynamic DNS tool base on DNSPod API"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		server.Command(),
		client.Command(),
	}
	app.Run(os.Args)
}
