package main

import (
	"github.com/urfave/cli"
	"net/http"
	"strings"
	"github.com/pantsing/log"
)

type Server struct {
	Listen   string
	CertFile string
	KeyFile  string
}

func (s *Server) Command() cli.Command {
	var c cli.Command
	c.Name = "server"
	c.Usage = "Return client internet IP"
	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf,c",
			Usage: "The configure file",
			Value: "gohome.toml",
		},
	}
	c.Action = s.Start
	return c
}

func (s *Server) Start(ctx *cli.Context) (err error) {
	err = loadConfigure(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	http.HandleFunc("/ip", s.IPHandler)
	err = http.ListenAndServeTLS(s.Listen, s.CertFile, s.KeyFile, nil)
	return err
}

func (s *Server) IPHandler(rw http.ResponseWriter, req *http.Request) {
	ip := strings.SplitN(req.RemoteAddr, ":", 2)[0]
	rw.Header().Set("Content-Type", "text/plain")
	rw.Write([]byte(ip))
}
