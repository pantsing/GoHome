package main

import (
	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
	"io/ioutil"
	"path/filepath"
)

type Conf struct {
	Server *Server
	Client *Client
}

type DDNSConf struct {
	DNSPOD struct {
		LoginToken string
		Format     string
	}
}

func loadConfigure(ctx *cli.Context) (err error) {
	filename, err := filepath.Abs(ctx.String("conf"))
	if err != nil {
		return
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = toml.Unmarshal(b, conf)
	return
}
