package main

import (
	"encoding/xml"
	"log"

	"github.com/foomo/soap"
)
//Cliente usando "github.com/foomo/soap"
// Request a simple request
type Request struct {
	XMLName xml.Name `xml:"request"`
	IDPerson	string `xml:"person"`
}

// Response a simple response
type Response struct {
	XMLName xml.Name `xml:"personResponse"`
	ID 		string   `xml:"id"`
	Name	string   `xml:"name"`
	Cedula	string   `xml:"cedula"`
}

func main() {
	soap.Verbose = true
	
	client := soap.NewClient("http://127.0.0.1:8090/person", nil, nil)
	
	response := &Response{}
	
	httpResponse, err := client.Call("GetPerson", Request{ IDPerson: "01",}, response)
	
	log.Println("\n\n")
	if err != nil {
		panic(err)
	}
	log.Println(response.Name, httpResponse.Status)
}
