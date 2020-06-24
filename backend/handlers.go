package main

import (
	"database/sql"
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

type ItemsList struct {
	Items []DomainSearched `json:"items"`
}
type DomainSearched struct {
	Domain string `json:"domain"`
}

func GetQueriedDomains(ctx *fasthttp.RequestCtx) {
	dbconnectorP, err := NewDBConnector("endpoints_admin", "defaultdb", "26257")
	if err != nil {
		fmt.Println("Hubieron problemas creando el conector a la BD")
	}
	rows, errQuery := dbconnectorP.connection.Query("SELECT dominio FROM endpoint_table GROUP BY dominio;")
	if errQuery != nil {
		log.Fatal(errQuery)
	}
	defer rows.Close()
	listDomains := make([]DomainSearched, 0)
	for rows.Next() {
		var dominio string
		errScan := rows.Scan(&dominio)
		if errScan != nil {
			log.Fatal(errScan)
		}
		temporalStruct := &DomainSearched{
			Domain: dominio,
		}
		listDomains = append(listDomains, *temporalStruct)
	}
	itemList := &ItemsList{
		Items: listDomains,
	}
	itemListMarshall, errMarshall := json.Marshal(itemList)
	if errMarshall != nil {
		log.Fatal(errMarshall)
	}
	ctx.WriteString(string(itemListMarshall))
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
	hrefLogo := GetHrefLinkLogo(domain)
	dbconnectorP, err3 := NewDBConnector("endpoints_admin", "defaultdb", "26257")
	if err3 != nil {
		fmt.Println("Hubieron problemas creando el conector a la BD")
	} else {
		fmt.Println("Conexión establecida!")
	}
	var dominio string
	var serversChanged bool
	var previousSslGrade string
	var jsonFromServerDB []byte
	var jsonUnmarshalled []ServersDef
	valueOfJSONByte, _ := json.Marshal(serversDefInstantiation)
	valueOfJSONString := string(valueOfJSONByte)
	errQuery := dbconnectorP.connection.QueryRow("SELECT dominio FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&dominio)
	if errQuery == sql.ErrNoRows {
		serversChanged = false
		previousSslGrade = "null"
		dbconnectorP.connection.Query("INSERT INTO defaultdb.endpoint_table VALUES ('" + domain + "','" + valueOfJSONString + "', '" + lowerGrade + "', now() AT TIME ZONE 'America/Bogota');")
	} else {
		errQuery2 := dbconnectorP.connection.QueryRow("SELECT grado_ssl FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&previousSslGrade)
		errQuery3 := dbconnectorP.connection.QueryRow("SELECT info_servers FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&jsonFromServerDB)
		if errQuery2 != nil {
			fmt.Println("Hubo un error obteniendo el grado anterior SSL de la DB")
		}
		if errQuery3 != nil {
			fmt.Println("Hubo un error obteniendo el JSON de servers de la DB")
		}
		serversChanged = false
		json.Unmarshal(jsonFromServerDB, &jsonUnmarshalled)
		for i := 0; i < len(serversDefInstantiation); i++ {
			if serversDefInstantiation[i] != jsonUnmarshalled[i] {
				serversChanged = true
			}
		}
	}
	jsonp := &JSONDef{
		Servers:          serversDefInstantiation,
		SslGrade:         lowerGrade,
		IsDown:           isServerDown,
		Title:            htmltitle,
		Logo:             hrefLogo,
		ServersChanged:   serversChanged,
		PreviousSslGrade: previousSslGrade,
	}
	jsonPrincipalByte, err2 := json.Marshal(jsonp)
	if err2 != nil {
		log.Fatal(err2)
	}
	ctx.WriteString(string(jsonPrincipalByte))
}
