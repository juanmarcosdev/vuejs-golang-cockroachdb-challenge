package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

// Este archivo consiste en los métodos que se usan para el Endpoint 1:
// * Obtener la IP de todos los servers/endpoints de un dominio y sus respectivos grados SSL (haciendo GET a la API de SSL Labs)
// * Obtener el país y la organización (owner) de cada server de un dominio (ejecutando el comando whois del sistema operativo)

// Este struct de tipo Domain tiene el host (el dominio, como google.com) y
// una slice de tipos Endpoint, que corresponden a los endpoints / servidores del dominio
// brindados por SSL Labs
type Domain struct {
	Host      string     `json:"host,omitempty"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

// Este struct de tipo Endpoint contiene la información solamente de la dirección IP
// y grado SSL, brindadas por la API de SSL Labs
type Endpoint struct {
	IPAddress string `json:"ipAddress,omitempty"`
	Grade     string `json:"grade,omitempty"`
}

// El struct de tipo MethodCaller es quien tendrá a disposición las funciones que se llaman
type MethodCaller struct{}

// Esta función (receiver del tipo MethodCaller) recibe una IP y ejecutando los comandos de whois
// en el sistema operativo logra retornar el país y owner de una IP en particular
func (mc *MethodCaller) ObtainIPCountryAndOrganization(ip string) (string, string) {
	// Este comado nos brinda el país de la IP
	commandS1 := []string{"-c", "whois " + ip + " | grep -i 'country' | cut -d ':' -f 2"}
	// Lo ejecutamos
	cmd1 := exec.Command("/bin/bash", commandS1...)
	// Este comando nos brinda el owner/organization de la IP
	commandS2 := []string{"-c", "whois " + ip + " | grep -i 'orgname' | cut -d ':' -f 2"}
	// Lo ejecutamos
	cmd2 := exec.Command("/bin/bash", commandS2...)
	// Obtenemos la standard output y el standard error de los comandos
	output1, err1 := cmd1.CombinedOutput()
	output2, err2 := cmd2.CombinedOutput()
	// Si hubo error en alguno de los comandos se notifica
	if err1 != nil {
		fmt.Printf("Hubo un error obteniendo el país del server")
	}
	if err2 != nil {
		fmt.Printf("Hubo un error obteniendo la organización del server")
	}
	// Ya habiendo obtenido las salidas, las  parseamos a string
	output1Str := string(output1)
	output2Str := string(output2)
	// Es importante eliminar los espacios a los que haya lugar en los strings, para tener la info limpia
	// y los retornamos
	return strings.TrimSpace(output1Str), strings.TrimSpace(output2Str)

}

// Esta función del struct MethodCaller recibe el dominio (como google.com)
// y devolverá la lista (en un slice) de tipos Endpoint que tenga la info de
// direcciones IP y grados SSL de cada endpoint/server de un dominio en particular
func (mc *MethodCaller) GetSSLLabsInfo(domain string) []Endpoint {
	// Creamos nuestro propio cliente http de GO para hacer la solicitud
	var clienteHTTP = &http.Client{Timeout: 20 * time.Second}
	// Hacemos el GET a la API de SSL Labs
	resp, err := clienteHTTP.Get("https://api.ssllabs.com/api/v3/analyze?host=" + domain)
	// Si hubieron errores en la solicitud, lo notificamos
	if err != nil {
		log.Fatal(err)
	}
	// Vamos a necesitar decodificar la salida JSON de la API externa
	decoder := json.NewDecoder(resp.Body)
	// Creamos struct de tipo Domain para allí guardar la info
	dom := &Domain{}
	// Decodificamos y filtramos los campos JSON que necesitamos en el tipo Domain
	errdef := decoder.Decode(dom)
	// Si hay errores los notificamos
	if errdef != nil {
		log.Fatal(errdef)
	}
	// Retornamos los Endpoints del dominio en el struct
	return dom.Endpoints
}
