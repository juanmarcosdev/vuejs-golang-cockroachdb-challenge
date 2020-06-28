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
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
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

// Esta función corresponde al primer Endpoint, donde se envía un dominio (POST)
// y se obtiene información sobre sus servidores.
func PostDomainAndGetInfo(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
	// El dominio a consultar es obtenido de la URL del endpoint
	domain, err := ctx.UserValue("domain").(string)
	// Si no se pudo obtener el dominio de la URL, se notifica el error (bad request)
	if !err {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
	}
	// Se instancia un struct del tipo MethodCaller (algo así como quien llama
	// los metodos, struct que tiene las funciones que se desarrollan en el endpoint)
	mc := &MethodCaller{}
	// Con el struct instanciado, llamamos a la función GetSSLLabsInfo,
	// que recibe el dominio que se quiere consultar, y devuelve un slice
	// de tipos de dato Endpoints, que llevan la información que se debe
	// extraer de la API de SSL Labs (cada endpoint con su ip y su grado SSL)
	endpoints := mc.GetSSLLabsInfo(domain)
	// Si de repente la solicitud a la API externa es muy rápida o por
	// alguna razón SSL Labs aún está obteniendo la info de el dominio,
	// se realiza un ciclo hasta que se obtenga la información.
	for len(endpoints) == 0 {
		time.Sleep(100 * time.Millisecond)
		endpoints = mc.GetSSLLabsInfo(domain)
	}
	// Se instancia un slice de tipos de dato ServersDef, con longitud
	// de la cantidad de endpoints, pues en cada uno de estos es donde
	// se va a compactar la información a mostrar
	serversDefInstantiation := make([]ServersDef, len(endpoints))
	// Se instancia la variable del país a la que pertenece un endpoint del dominio
	var country string
	// Igual para la organización (owner) del cada endpoint
	var organization string
	// Se realiza un ciclo recorriendo la cantidad de endpoints que haya,
	// y por cada uno de ellos se llama a la función ObtainIPCountryAndOrganization
	// que hace uso del comando whois, instalado en el contenedor base del Backend,
	// dicho comando nos provee la información de país & owner para cada IP de cada endpoint
	for i, e := range endpoints {
		// Se obtiene la data de la función
		country, organization = mc.ObtainIPCountryAndOrganization(e.IPAddress)
		// Se guarda en cada endpoint respectivamente y de una vez
		// se guarda la información de IP Address y SslGrade brindada anteriormente
		// y así ya se va construyendo la primera parte de Servers del JSON
		serversDefInstantiation[i] = ServersDef{
			Address:  e.IPAddress,
			SslGrade: e.Grade,
			Country:  country,
			Owner:    organization,
		}
	}
	// Acá se instancia un slice de los grados de cada endpoint, con el objetivo
	// de encontrar el menor.
	sliceOfGrades := make([]string, len(endpoints))
	// Se guardan los grados obtenidos
	for i, e := range endpoints {
		sliceOfGrades[i] = e.Grade
	}
	// Se procede a ordenarlos lexicográficamente, así, el más bajo
	// estará en la última posición (excepto en el caso en el que haya solo A y A+)
	sort.Strings(sliceOfGrades)
	// Aquí guardaremos el grado más bajo
	var lowerGrade string
	// Para el caso especial donde hayan solo servidores de grados A y A+,
	// al ordenarse A+ estará en la última posición, sin embargo esto es incorrecto,
	// por ende, si existe A en el slice, ya demás el último es A+, se corrige
	// y se pone como el grado más bajo el A.
	// Si esto no sucede, simplemente no aplica a este filtro y sigue el programa normal.
	isAinSlice := stringInSliceOfStrings("A", sliceOfGrades)
	if isAinSlice && sliceOfGrades[len(endpoints)-1] == "A+" {
		lowerGrade = "A"
	}
	// De no calificar en el caso de arriba simplemente se elige el de la última posición
	lowerGrade = sliceOfGrades[len(endpoints)-1]
	// Para poder saber si el servidor está down o no, se crea un pequeño cliente HTTP
	// con timeout de 20 segundos esperando respuesta. Además se aprovecha
	// y con esta respuesta GET se obtendrán los valores del TITLE y LOGO de la página.
	var clienteHTTP = &http.Client{Timeout: 20 * time.Second}
	// Se crea variable para almacenar si el server está down o no
	var isServerDown bool
	// Se obtiene la respuesta del servidor
	respo, err2 := clienteHTTP.Get("http://www." + domain)
	// De haber un error significa el que el dominio está caído
	if err2 != nil {
		log.Fatal(err2)
		isServerDown = true
	}
	// Lo último que hacemos es cerrar la respuesta
	defer respo.Body.Close()
	// Como no hubo error, el servidor no está down
	isServerDown = false
	// Aquí le pasamos el body de la respuesta a la función ObtenerHTMLTitle
	// que recibe dicho cuerpo y busca recursivamente en toda la página el
	// tag <title> del HTML y devuelve lo que hay adentro
	htmltitle, ok := ObtenerHTMLTitle(respo.Body)
	// Si no se pudo obtener el title, se notifica
	if ok == false {
		fmt.Println("Hubo un error obteniendo el title")
	}
	// Luego obtendremos el href (link) del logo del sitio.
	// Hay que tener en cuenta que solo funciona para sitios que
	// tengan su logo de esta manera:
	// html -> head -> link y en la etiqueta link tenga un atributo
	// "rel" con valor de "shortcut icon". Hará match con el que tenga
	// estas características y obtendrá el valor que está en el atributo "href".
	hrefLogo := GetHrefLinkLogo(domain)
	// Empezamos a crear el conector a nuestra DB para empezar a guardar la información
	dbconnectorP, err3 := NewDBConnector("endpoints_admin", "defaultdb", "26257")
	// Si hay problemas creando el conector DB se notifican
	if err3 != nil {
		fmt.Println("Hubieron problemas creando el conector a la BD")
	}
	// Variable dominio que sirve para extraer el campo "dominio" de la tabla en la DB
	var dominio string
	// Aquí sabremos si los servidores han cambiado entre un registro y una nueva consulta
	var serversChanged bool
	// Aquí obtendremos el grado SSL de la consulta anterior:
	// Si es primera vez que el dominio se consulta, este campo tendrá
	// un valor de "null". Ya a la segunda vez que pueda consultar el primer registro (de una hora o antes)
	// si le pondrá el valor respectivo.
	var previousSslGrade string
	// En este slice de bytes obtendremos la información correspondiente a los servidores del dominio
	// guardados en la BD
	var jsonFromServerDB []byte
	// Luego el JSON de la DB lo transformaremos a slice de ServersDef para
	// poder compararlo con el que se consulta in-time
	var jsonUnmarshalled []ServersDef
	// El que obtuvimos haciendo las consultas a SSL Labs y whois, lo
	// casteamos slice de bytes y luego a String, para poder enviarlo a la DB
	valueOfJSONByte, _ := json.Marshal(serversDefInstantiation)
	valueOfJSONString := string(valueOfJSONByte)
	// Realizamos una consulta para saber si del dominio consultado existen registros
	// de hace una hora o más antes (59 minutos, 40 minutos, 30 minutos, 1 minuto).
	// Si existen registros que han permanecido por 1 hora y 1 segundo, los ignora y crea
	// un nuevo registro. Si existen registros desde 1 segundo hasta 1 hora, los tiene en cuenta
	// y con ellos obtiene la info de servers_changed y previous_ssl_grade.
	// (También de una vez obtiene el valor de dominio de la BD)
	errQuery := dbconnectorP.connection.QueryRow("SELECT dominio FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&dominio)
	if errQuery == sql.ErrNoRows {
		// Si la consulta no arroja resultados, no podemos decir que los servidores hayan cambiado
		// pues no hay punto de comparación, entonces dicho atributo se setea en false.
		serversChanged = false
		// Como no hay registro de hace una hora o más antes, se pone en "null" el campo de previous_ssl_grade
		previousSslGrade = "null"
		// Como no hay registro de hace una hora o más antes, la información que hemos consultado in-time se almacena en la BD para futuras consultas.
		dbconnectorP.connection.Query("INSERT INTO defaultdb.endpoint_table VALUES ('" + domain + "','" + valueOfJSONString + "', '" + lowerGrade + "', now() AT TIME ZONE 'America/Bogota');")
	}
	// Si existe registro de hace una hora o más antes, vamos a obtener el grado SSL anterior en dicha consulta y la información de los servers.
	errQuery2 := dbconnectorP.connection.QueryRow("SELECT grado_ssl FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&previousSslGrade)
	errQuery3 := dbconnectorP.connection.QueryRow("SELECT info_servers FROM endpoint_table WHERE dominio = '" + domain + "' AND hora_consulta > NOW() AT TIME ZONE 'America/Bogota' - INTERVAL '1 hour';").Scan(&jsonFromServerDB)
	// Si alguna de las consultas falló, notificarlo
	if errQuery2 != nil {
		fmt.Println("Hubo un error obteniendo el grado anterior SSL de la DB")
	}
	if errQuery3 != nil {
		fmt.Println("Hubo un error obteniendo el JSON de servers de la DB")
	}
	// En inicios se supone que los servidores no han cambiado
	serversChanged = false
	// Aquí vamos a formatear a JSON los datos obtenidos de servers de la DB
	// (Con el formato de []ServersDef)
	json.Unmarshal(jsonFromServerDB, &jsonUnmarshalled)
	// Ahora aquí tenemos dos slices de ServersDef. Lo que haremos será comparar elemento a elemento,
	// es decir, endpoint a endpoint, al ser structs pueden compararse entre ellos y si difieren en algo pequeño,
	// serán structs distintos y servers_changed será true. Si resultan ser todo lo mismo
	// quiere decir que los servers no han cambiado y no se cambiará su valor al iniciado arriba (Seguirá siendo false que hayan cambiado)
	for i := 0; i < len(serversDefInstantiation); i++ {
		if serversDefInstantiation[i] != jsonUnmarshalled[i] {
			serversChanged = true
		}
	}

	// Finalmente, instanciamos un struct del tipo JSONDef, que nos servirá
	// para mostrar toda la información en el navegador o como respuesta
	// a la petición POST. En cada campo le asignamos la información obtenida anteriormente
	jsonp := &JSONDef{
		Servers:          serversDefInstantiation,
		SslGrade:         lowerGrade,
		IsDown:           isServerDown,
		Title:            htmltitle,
		Logo:             hrefLogo,
		ServersChanged:   serversChanged,
		PreviousSslGrade: previousSslGrade,
	}
	// Luego vamos a pasar de struct -> JSON (en slice de bytes)
	jsonPrincipalByte, err2 := json.Marshal(jsonp)
	// Si hubo un error en el parseo de struct a JSON, notificar
	if err2 != nil {
		log.Fatal(err2)
	}
	// Finalmente vamos a mostrar como respuesta el JSON de bytes casteado
	// a string y nuestro endpoint a finalizado.
	ctx.WriteString(string(jsonPrincipalByte))
}
