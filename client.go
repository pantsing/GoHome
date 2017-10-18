package main

import (
	"crypto/tls"
	"github.com/pantsing/log"
	"github.com/urfave/cli"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Client struct {
	currentIP          net.IP
	ddns               map[string]DDNS
	CheckInterval      int // check IP interval seconds
	ServerAPI          string
	InsecureSkipVerify bool
	DDNSConf           DDNSConf
	Domains            map[string]struct {
		DDNS    string
		Records []Record
	}
}

type Record struct {
	Host       string //主机记录   可选     (如:www,默认为@)
	RecordType string //记录类型   必选     (通过API记录类型获得，大写英文，如:A)
	RecordLine string //记录线路   必选     (通过API记录线路获得，中文，比如：默认)
	//Value      string //记录值     必选     (如 IP:200.200.200.200, CNAME: cname.dnspod.com., MX: mail.dnspod.com., )
	//MX         int    //MX优先级   MX记录必选     (当记录类型是 MX 时有效，范围1-20)
	//TTL        int    //TTL       可选     (范围1-604800，不同等级域名最小值不同)
}

func (client *Client) Command() cli.Command {
	var c cli.Command
	c.Name = "client"
	c.Usage = "Get Internet IP and Update DNSPod"
	c.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "conf,c",
			Usage: "The configure file",
			Value: "gohome.toml",
		},
	}
	c.Action = client.Start
	return c
}

func (client *Client) Start(ctx *cli.Context) (err error) {
	err = loadConfigure(ctx)
	if err != nil {
		log.Error(err)
		return
	}

	client.ddns = map[string]DDNS{"DNSPOD": NewDNSPodAPI(client.DDNSConf)}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: client.InsecureSkipVerify},
	}
	httpClient := &http.Client{Transport: tr}
	ticker := time.NewTicker(time.Duration(client.CheckInterval) * time.Second)
	for range ticker.C {
		resp, err := httpClient.Get(conf.Client.ServerAPI)
		if err != nil {
			log.Warn(err)
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Warn(err)
			continue
		}
		resp.Body.Close()
		ipStr := string(b)
		ip := net.ParseIP(ipStr)
		if ip != nil && client.currentIP.String() != ip.String() {
			client.updateRecords(ip)
		}
	}
	return nil
}

func (client *Client) updateRecords(newIP net.IP) {
	var err error
LOOP:
	for domain := range client.Domains {
		ddnsName := client.Domains[domain].DDNS
		ddns := client.ddns[ddnsName]
		for _, r := range client.Domains[domain].Records {
			err = ddns.Update(domain, client.currentIP, r)
			if err != nil {
				break LOOP
			}
		}
	}

	if err == nil {
		client.currentIP = newIP
	} else {
		log.Warn(err)
		time.Sleep(time.Minute)
	}
}
