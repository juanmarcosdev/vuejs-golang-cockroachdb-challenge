package main

import (
	"fmt"
	"os/exec"
	"strings"
)

type MethodCaller struct{}

func (md *MethodCaller) ObtainServerCountryAndOrganization(domain string) (string, string) {
	commandS1 := []string{"-c", "whois " + domain + " | grep -i 'registrant country' | cut -d ' ' -f 3"}
	cmd1 := exec.Command("/bin/bash", commandS1...)
	commandS2 := []string{"-c", "whois " + domain + " | grep -i 'registrant organization' | awk -F':' '{print $2}'"}
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
	return output1Str, strings.TrimSpace(output2Str)

}
