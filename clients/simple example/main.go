package main


import(
	"fmt"
	"net/http"
	"bytes"
	"encoding/xml"
	"text/template"
	"io/ioutil"
)

// Objeto que sera enviado en la petición con el id de la persona
type Request struct {
	XMLName xml.Name `xml:"request"`
	IDPerson     string `xml:"person"`
}

//Objeto que recibira la respuesta
type Response struct {
   //Envelope es la primera etiqueta del archivo XML que vmos a leer
   XMLName  xml.Name `xml:"Envelope"`
   SoapBody *SOAPBodyResponse
}

type SOAPBodyResponse struct {
   //Accedemos a la etiqueta BODY
   XMLName      xml.Name `xml:"Body"`
   //Fault es el nombre de la etiqueta que contien el error que devuelve el servicio, en caso de que alla ocurrido un inconveniente
   FaultDetails *Fault
   Resp 		    *PersonResponse 
}

type Fault struct {
   XMLName     xml.Name `xml:"Fault"`
   Faultcode   string   `xml:"faultcode"`
   Faultstring string   `xml:"faultstring"`
}

type PersonResponse struct {
   //Accedemos a la etiquieta personResponse que es donde esta la información que me devuelve el servicio
   XMLName       xml.Name `xml:"personResponse"`
	ID			     string   `xml:"id"`
   Name          string   `xml:"name"`
   Cedula        string   `xml:"cedula"`
} 

//Creamos la estructura de la petición que sera enviada al servicio
var getTemplate =`
<soapenv:Envelope
 xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/"
 xmlns:api="http://soapdummies.com/api">
 <soapenv:Header/>
 <soapenv:Body>
   	<request>
        <person>{{.IDPerson}}</person>
	   </request>
 </soapenv:Body>
</soapenv:Envelope>
`

const URL = "http://127.0.0.1:8090/person"

func main(){
	callSOAPClientSteps()
}

func callSOAPClientSteps(){
   //Generamos el cuerpo de la petición
   req := generateRequest()

   //Generamos la petición http que debe ser siempre del tipo Post
   httpReq, err := generateSOAPRequest(req)
   if err != nil {
      fmt.Println("Some problem occurred in request generation: ",err)
      return
   }

   //Llamamos al servicio
   response, err := soapCall(httpReq)
   if err != nil {
      fmt.Println("Problem occurred in making a SOAP call: ", err)
      return
   }

   fmt.Println("RESULTADO:")
   fmt.Println(response.SoapBody.Resp.ID)
   fmt.Println(response.SoapBody.Resp.Name)
   fmt.Println(response.SoapBody.Resp.Cedula)
}

func generateRequest() *Request {
   req := Request{}
   req.IDPerson = "01"
   return &req
}

func generateSOAPRequest(req *Request) (*http.Request, error) {
   // Using the var getTemplate to construct request
   template, err := template.New("InputRequest").Parse(getTemplate)
   if err != nil {
      fmt.Println("Error while marshling object. %s ",err.Error())
      return nil,err
   }

   doc := &bytes.Buffer{}
   // Replacing the doc from template with actual req values
   err = template.Execute(doc, req)
   if err != nil {
      fmt.Println("template.Execute error. %s ",err.Error())
      return nil,err
   }

   buffer := &bytes.Buffer{}
   encoder := xml.NewEncoder(buffer)
   err = encoder.Encode(doc.String())
   if err != nil {
      fmt.Println("encoder.Encode error. %s ",err.Error())
      return nil,err
   }

   r, err := http.NewRequest(http.MethodPost, URL,bytes.NewBuffer([]byte(doc.String())))
   if err != nil {
      fmt.Println("Error making a request. %s ", err.Error())
      return nil,err
   }

   r.Header.Set("Soapaction", "GetPerson")

   return r, nil
}

func soapCall(req *http.Request) (*Response, error) {
   client := &http.Client{}
   resp, err := client.Do(req)

   if err != nil{
      return nil,err
   }

   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      return nil,err
   }
   defer resp.Body.Close()

   r := &Response{}
   err = xml.Unmarshal(body, &r)

   if err != nil {
      return nil,err
   }


   return r, nil
}


/*
   XML Recibido
   <Envelope xmlns="http://www.w3.org/2003/05/soap-envelope">
        <Header xmlns="http://www.w3.org/2003/05/soap-envelope"></Header>
         <Body xmlns="http://www.w3.org/2003/05/soap-envelope">
               <personResponse>
                       <id>01</id>
                       <name>Manuel</name>
                       <cedula>24523126</cedula>
               </personResponse>
        </Body>
</Envelope>
*/
