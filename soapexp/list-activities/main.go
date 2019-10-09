package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"rits/carriage/api/models/frontierlink"
	"rits/carriage/api/soap"
	"text/template"
)

var templates *template.Template

type Body struct {
	ListActivities *frontierlink.ListActivities `xml:"ListActivitiesResponse,omitempty" json:"ListActivitiesResponse,omitempty"`
}

type Envelope struct {
	Body *Body `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body,omitempty" json:"Body,omitempty"`
}

func main() {
	url := "https://frontierlink-cert.aapt.com.au:8482/FrontierLink/services/ServiceStatus"
	action := ""

	templates, err := soap.ParseTemplates("../../templates")
	if err != nil {
		fmt.Println(err)
	}

	type ListActivitiesRequest struct {
		ProvisioningCaseID string
	}

	var data ListActivitiesRequest
	data = ListActivitiesRequest{ProvisioningCaseID: "CA00002"}

	var payload bytes.Buffer
	s4 := templates.Lookup("ListActivitiesRequest.gohtml")
	s4.ExecuteTemplate(&payload, "ListActivitiesRequest", data)

	var Result Envelope
	body, err := soap.FrontierLinkSoapHandler(url, action, payload.Bytes(), "../../certs/cert.pem")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("body response= .....")
	fmt.Println(string(body))

	err = xml.Unmarshal(body, &Result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("line 39 result=")
	fmt.Println(Result)
	spew.Dump(Result)
}
