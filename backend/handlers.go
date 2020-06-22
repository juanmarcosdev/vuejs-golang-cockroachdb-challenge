package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

func stringInSliceOfStrings(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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
}

func PostDomainAndGetInfo(ctx *fasthttp.RequestCtx) {
	domain, err := ctx.UserValue("domain").(string)
	if !err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	mc := &MethodCaller{}
	endpoints := mc.GetSSLLabsInfo(domain)
	serversDefInstantiation := make([]ServersDef, len(endpoints))
	var country string
	var organization string
	for i, e := range endpoints {
		country, organization = mc.ObtainIPCountryAndOrganization(e.IPAddress)
		serversDefInstantiation[i] = ServersDef{
			Address:  e.IPAddress,
			SslGrade: e.Grade,
			Country:  country,
			Owner:    organization,
		}
	}
	sliceOfGrades := make([]string, len(endpoints))
	for i, e := range endpoints {
		sliceOfGrades[i] = e.Grade
	}
	sort.Strings(sliceOfGrades)
	var lowerGrade string
	isAinSlice := stringInSliceOfStrings("A", sliceOfGrades)
	if isAinSlice && sliceOfGrades[len(endpoints)-1] == "A+" {
		lowerGrade = "A"
	}
	lowerGrade = sliceOfGrades[len(endpoints)-1]
	var clienteHTTP = &http.Client{Timeout: 20 * time.Second}
	var isServerDown bool
	respo, err2 := clienteHTTP.Get("http://www." + domain)
	if err2 != nil {
		log.Fatal(err2)
		isServerDown = true
	}
	defer respo.Body.Close()
	isServerDown = false
	htmltitle, ok := ObtenerHTMLTitle(respo.Body)
	if ok == false {
		fmt.Println("Hubo un error obteniendo el title")
	}
	jsonp := &JSONDef{
		Servers:  serversDefInstantiation,
		SslGrade: lowerGrade,
		IsDown:   isServerDown,
		Title:    htmltitle,
	}
	jsons, err2 := json.Marshal(jsonp)
	if err2 != nil {
		log.Fatal(err2)
	}
	ctx.WriteString(string(jsons))
}
