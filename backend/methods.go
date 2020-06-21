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

type Domain struct {
	Host      string     `json:"host,omitempty"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
}

type Endpoint struct {
	IPAddress string `json:"ipAddress,omitempty"`
	Grade     string `json:"grade,omitempty"`
}

type MethodCaller struct{}

func (mc *MethodCaller) ObtainIPCountryAndOrganization(ip string) (string, string) {
	commandS1 := []string{"-c", "whois " + ip + " | grep -i 'country' | cut -d ':' -f 2"}
	cmd1 := exec.Command("/bin/bash", commandS1...)
	commandS2 := []string{"-c", "whois " + ip + " | grep -i 'orgname' | cut -d ':' -f 2"}
	cmd2 := exec.Command("/bin/bash", commandS2...)
	output1, err1 := cmd1.CombinedOutput()
	output2, err2 := cmd2.CombinedOutput()
	if err1 != nil {
		fmt.Printf("Hubo un error obteniendo el país del server")
	}
	if err2 != nil {
		fmt.Printf("Hubo un error obteniendo la organización del server")
	}
	output1Str := string(output1)
	output2Str := string(output2)
	return strings.TrimSpace(output1Str), strings.TrimSpace(output2Str)

}

func (mc *MethodCaller) GetSSLLabsInfo(domain string) []Endpoint {
	var clienteHTTP = &http.Client{Timeout: 20 * time.Second}
	resp, err := clienteHTTP.Get("https://api.ssllabs.com/api/v3/analyze?host=" + domain)
	if err != nil {
		log.Fatal(err)
	}
	decoder := json.NewDecoder(resp.Body)
	dom := &Domain{}
	errdef := decoder.Decode(dom)
	if errdef != nil {
		log.Fatal(errdef)
	}
	return dom.Endpoints
}
