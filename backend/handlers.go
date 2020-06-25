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

// El archivo handlers.go es el más importante de la estructura del servidor,
// puesto que aquí están todos los elementos que ayudan a hacer la tarea
// de los dos endpoints asociados a los handlers aquí escritos.

// Esta función recibe un string a, un slice de strings 'list', y retorna un
// booleano si el string a pertenece al slice de strings 'list'.
// Esta función se utiliza para realizar una ponderación correcta
// de los grados SSL de un servidor, puesto que el campo ssl_grade (del JSON
// que retorna el Endpoint que registra la info de los servers) lleva
// el grado SSL más bajo (en la escala dada por SSL Labs). Para poder obtenerlo
// se ordena de manera lexicográfica el slice que contiene todos los grados
// de cada servidor y se escoge el de la última posición (que sería el grado más bajo).
// Sin embargo, existe un caso donde solo hayan servidores con grados A y A+.
// Según SSL Labs, A es menor que A+. Sin embargo, al realizar el ordenado
// quedará de último A+, por orden lexicográfico. Esta función ayuda a saber si
// el último resulta ser A+ y existe un grado A, para poder imprimir el grado A
// que sería la respuesta correcta.
func stringInSliceOfStrings(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Este struct JSONDef es el struct que es homólogo en forma al JSON que
// debe retornar el primer endpoint (POST) de escritura de información.
// Tiene exactamente todos los campos del JSON que debe retornar, y en particular
// tiene un campo que es un slice de tipos struct ServersDef.
type JSONDef struct {
	Servers          []ServersDef `json:"servers"`
	ServersChanged   bool         `json:"servers_changed"`
	SslGrade         string       `json:"ssl_grade"`
	PreviousSslGrade string       `json:"previous_ssl_grade"`
	Logo             string       `json:"logo"`
	Title            string       `json:"title"`
	IsDown           bool         `json:"is_down"`
}

// Este struct es la información correspondiente a los servidores de un
// dominio respectivo (esta información la brinda la API de SSL Labs). Contiene
// todos los campos homólogos a lo que debe haber en el JSON-array de "servers"
// en el primer Endpoint.
type ServersDef struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
	Country  string `json:"country"`
	Owner    string `json:"owner"`
}

// Este struct ItemsList corresponde al segundo endpoint,
// puesto que tiene un elemento "items" que es un slice de tipo DomainSearched,
// es decir que este struct es homólogo a lo que retorna en JSON el segundo endpoint de consulta (GET).
type ItemsList struct {
	Items []DomainSearched `json:"items"`
}

// Este tipo DomainSearched corresponde a un JSON sencillo, con un solo campo
// "domain", donde debe ir registrado un dominio que haya sido consultado
// (que se encuentre en la base de datos).
type DomainSearched struct {
	Domain string `json:"domain"`
}

// Este es el handler correspondiente al Endpoint 2, donde se pueden obtener
// los dominios consultados.
func GetQueriedDomains(ctx *fasthttp.RequestCtx) {
	// Lo primero que se hace es obtener una conexión a la base de datos de CockroachDB que está en el otro contenedor de Docker
	dbconnectorP, err := NewDBConnector("endpoints_admin", "defaultdb", "26257")
	// Si hubieron errores creando la conexión se notifican
	if err != nil {
		fmt.Println("Hubieron problemas creando el conector a la BD")
	}
	// Aquí se hace la consulta a la base de datos, sobre todos los dominios consultados sin importar hace cuánto se hicieron
	// La cláusula GROUP BY en el Query es fundamental, puesto que si se consultaron varias veces un mismo dominio
	// no obtendremos repetidos.
	rows, errQuery := dbconnectorP.connection.Query("SELECT dominio FROM endpoint_table GROUP BY dominio;")
	// Si hubieron errores en la consulta, notificarlos
	if errQuery != nil {
		log.Fatal(errQuery)
	}
	// Al final del programa debemos de cerrar las filas que arrojó el Query
	defer rows.Close()
	// Instanciamos un slice que contendrá structs del tipo DomainSearched,
	// aquí se van a guardar todos los dominios que arroje la consulta
	listDomains := make([]DomainSearched, 0)
	// Recorremos todas las filas que se obtuvieron del Query
	for rows.Next() {
		var dominio string
		// Por cada fila vamos a extraer (Scan) el campo de nombre "dominio"
		errScan := rows.Scan(&dominio)
		// Si hay un error obteniendo el campo se notifica
		if errScan != nil {
			log.Fatal(errScan)
		}
		// Habiendo guardado el valor de "dominio" en la variable con el mismo nombre,
		// creamos un struct temporal que tiene dicho campo y que siempre se va reemplazando, pero que su
		// información persiste gracias a que luego se guardan en el slice listDomains
		temporalStruct := &DomainSearched{
			Domain: dominio, // Aquí le asignamos al struct temporal dicho valor
		}
		// Y para conservarlo lo guardamos en listDomains, al ser declarada fuera del scope del for
		// se conservará
		listDomains = append(listDomains, *temporalStruct)
	}
	// Tras haber obtenido todos los dominios y haberlos guardado en la lista,
	// creamos el struct itemList que va a tener como atributo la lista que
	// acabamos de construir
	itemList := &ItemsList{
		Items: listDomains,
	}
	// Vamos a convertir dicha lista a un slice de bytes que tendrán formato JSON
	itemListMarshall, errMarshall := json.Marshal(itemList)
	// Si hubo un error en el parseo de struct a JSON lo notificamos
	if errMarshall != nil {
		log.Fatal(errMarshall)
	}
	// Finalmente escribiremos en la respuesta del navegador el contenido
	// del JSON convertido a string
	ctx.WriteString(string(itemListMarshall))
}

func PostDomainAndGetInfo(ctx *fasthttp.RequestCtx) {
	domain, err := ctx.UserValue("domain").(string)
	if !err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	mc := &MethodCaller{}
	endpoints := mc.GetSSLLabsInfo(domain)
	for len(endpoints) == 0 {
		time.Sleep(100 * time.Millisecond)
		endpoints = mc.GetSSLLabsInfo(domain)
	}
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
