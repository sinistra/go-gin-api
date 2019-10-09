package soap

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
	"time"
)

func FrontierLinkSoapHandler(url string, action string, payload []byte, certFile string) ([]byte, error) {

	timeout := time.Duration(30 * time.Second)
	cert, err := tls.LoadX509KeyPair(certFile, certFile)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//RootCAs: caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/xml, multipart/related")
	req.Header.Set("SOAPAction", action)
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	return bodyBytes, nil
}

func ParseTemplates(path string) (*template.Template, error) {
	var allFiles []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalln(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".gohtml") {
			allFiles = append(allFiles, path+"/"+filename)
		}
	}

	//parses all .gohtml files in the 'templates' folder
	return template.ParseFiles(allFiles...)
}
