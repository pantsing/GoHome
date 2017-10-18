package main

import (
	"encoding/json"
	"fmt"
	"github.com/pantsing/log"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

type DDNS interface {
	Update(domain string, ip net.IP, r Record) error
}

type DNSPodAPI struct {
	token  string
	format string
}

func NewDNSPodAPI(conf DDNSConf) *DNSPodAPI {
	return &DNSPodAPI{
		token:  conf.DNSPOD.LoginToken,
		format: conf.DNSPOD.Format,
	}
}

type DNSPodAPIResp struct {
	Status struct {
		Code       string
		Message    string
		Created_at string
	}
	Domain struct {
		ID         string
		Name       string
		punycode   string
		grade      string
		Owner      string
		Ext_status string
		TTL        int
		Min_TTL    int
	}
	Info struct {
		Sub_domains  string
		Record_total string
	}
	Records []struct {
		ID             string
		Name           string
		Value          string
		Line           string
		Line_ID        string
		Type           string
		TTL            string
		Weight         string
		MX             string
		Enable         string
		Status         string
		Monitor_status string
		Remark         string
		Updated_on     string
		Use_aqb        string
	}
}

func (api *DNSPodAPI) Update(domain string, ip net.IP, r Record) (err error) {
	values := url.Values{}
	values.Set("login_token", api.token)
	values.Set("format", api.format)
	values.Set("domain", domain)
	values.Set("sub_domain", r.Host)

	var resp *http.Response
	resp, err = http.PostForm("https://dnsapi.cn/Record.List", values)
	if err != nil {
		log.Warn(err)
		return
	}
	var b []byte
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(err)
		return
	}
	resp.Body.Close()
	data := new(DNSPodAPIResp)
	err = json.Unmarshal(b, data)
	if err != nil {
		log.Warn(err, string(b))
		return
	}
	if data.Status.Code != "1" {
		err = fmt.Errorf("%s %s %v", data.Status.Code, data.Status.Message, data.Records)
		log.Warn(err)
		return
	}

	if len(data.Records) == 0 {
		log.Warnf("Record not found: %v", values)
		return
	}

	values.Set("record_id", data.Records[0].ID)
	values.Set("record_line", r.RecordLine)
	values.Set("value", ip.String())
	log.Info(values)

	resp, err = http.PostForm("https://dnsapi.cn/Record.Ddns", values)
	if err != nil {
		log.Warn(err)
		return
	}
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Warn(err)
		return
	}
	resp.Body.Close()
	log.Info(string(b))

	data = new(DNSPodAPIResp)
	err = json.Unmarshal(b, data)
	if err != nil {
		log.Warn(err, string(b))
		return
	}
	if data.Status.Code != "1" {
		err = fmt.Errorf("%s %s %v", data.Status.Code, data.Status.Message, data.Records)
	}
	return
}
