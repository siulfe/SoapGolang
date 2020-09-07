package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/foomo/soap"
)
//Servicio usando "github.com/foomo/soap"

// Objeto que llegara en las peticiones que hagan a el servicio
type Request struct {
	//La información estara contenida dentro de la etiqueta "request"
	XMLName xml.Name `xml:"request"`
	//El id de la persona esta en la etiqueta "person"
	IDPerson	string `xml:"person"`
}

// Objeto que devolveremos en nuestras respuestas
type Response struct {
	XMLName xml.Name `xml:"personResponse"`
	ID 		string   `xml:"id"`
	Name	string   `xml:"name"`
	Cedula	string   `xml:"cedula"`
}

// Objeta para devolver un listado de Personas
type Responses struct {
	XMLName xml.Name `xml:"personsResponses"`
	Persons []Response
}

func main(){
	server := soap.NewServer()
	server.UseSoap12()
	server.RegisterHandler(
		// Ruta
		"/person",
		// SOAPAction - Identificador que debe llamar el cliente cuando hace la petición
		"GetPerson",
		// tagname of soap body content
		"request",
		// RequestFactoryFunc - give the server sth. to unmarshal the request into. Le pasamos como una función el objeto que queremos pasarle al operationHandlerFunc
		func() interface{} {
			return &Request{}
		},
		// OperationHandlerFunc - en donde procesamos la petición y damos una respuesta
		handlerPerson,
	)

	server.RegisterHandler("/persons", "GetPersons", "request",func() interface{} {
			return &Request{}
		}, handlerPersons)
	
	//Subimos el servidor
	err := server.ListenAndServe(":8090")

	if err != nil{
		fmt.Println("Error al subir servidor: ",err)
	}

}


func handlerPerson(request interface{}, w http.ResponseWriter, httpRequest *http.Request) (response interface{}, err error) {
	fmt.Println("request: ", request)
	//Transformamos la interface en un tipo Request
	resp := request.(*Request)
	//Obtenemos la persona
	response = getPerson(resp.IDPerson)

	return
}

func handlerPersons(request interface{}, w http.ResponseWriter, httpRequest *http.Request) (response interface{}, err error) {
	fmt.Println("handler Persons active")
	
	response = Responses{
		//Obtenermos el listado de personas
		Persons: getPersons(),
	}
	return
}


func getPerson(id string) *Response{
	resp := &Response{
		ID: id,
		Name: "Manuel",
		Cedula: "24523126",
	}

	return resp
}


func getPersons() []Response{
	resp := []Response{}

	for i := 0; i < 5; i++ {
		person := Response{
			ID: fmt.Sprintf("%d",i),
			Name: "Manuel"+fmt.Sprintf("%d",i),
			Cedula: "2452312"+fmt.Sprintf("%d",i),
		}
		resp = append(resp,person)
	}

	return resp
}



