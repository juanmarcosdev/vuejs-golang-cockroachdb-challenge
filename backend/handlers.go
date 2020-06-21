package main

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

type JSONDef struct {
	Servers          []ServersDef `json:"servers"`
	ServersChanged   bool         `json:"servers_changed"`
	SslGrade         string       `json:"ssl_grade"`
	PreviousSslGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
}

type ServersDef struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

func GetQueriedDomains(ctx *fasthttp.RequestCtx) {
	ctx.WriteString("Get funcionando!")
	// mc := &MethodCaller{}
	// endpoints := mc.GetSSLLabsInfo("netflix.com")
}

func PostDomainAndGetInfo(ctx *fasthttp.RequestCtx) {
	domain, err := ctx.UserValue("domain").(string)
	if !err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	mc := &MethodCaller{}
	endpoints := mc.GetSSLLabsInfo(domain)
	servers_def := make([]ServersDef, len(endpoints))
	var country string
	var organization string
	for i, e := range endpoints {
		country, organization = mc.ObtainIPCountryAndOrganization(e.IPAddress)
		servers_def[i] = ServersDef{
			Address:  e.IPAddress,
			SslGrade: e.Grade,
			Country:  country,
			Owner:    organization,
		}
	}
	jsonp := &JSONDef{
		Servers: servers_def,
	}
	jsons, err2 := json.Marshal(jsonp)
	if err2 != nil {
		log.Fatal(err2)
	}
	ctx.WriteString(string(jsons))
	// country, organization := mc.ObtainIPCountryAndOrganization(domain)
	// ctx.WriteString(country + "\n" + organization)
}
