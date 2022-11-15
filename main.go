package main

import (
	"fmt"
	"log"
	"net/http"
)

var key1 = "KEY1"
var key2 = "KEY2"

type Decider struct {
	service map[string]Service
}

type Service struct {
	ClientName string
	Parser     IParser
}

type IParser interface {
	Parse(in string) (out string) // TODO: change to Parser Signature
}

func (s *Service) Serve(in string) {
	result := s.Parser.Parse(in)
	fmt.Println("parsed result: ", result)
}

func (d *Decider) AddCustomer(name string, service Service) {
	// Initiate map to avoid nil panic
	if d.service == nil {
		d.service = make(map[string]Service)
	}
	d.service[name] = service
}

func (d *Decider) Decide(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Start Deciding.")
		var realTimeClient = "kun"
		var ToFull = make(map[string]string)
		ToFull["clientName"] = realTimeClient
		for _, toAddress := range ToFull {
			if service, ok := d.service[toAddress]; ok {
				service.Serve("Hi")
			}
		}

		next.ServeHTTP(w, r)
		log.Print("End Deciding.")
	})
}

type KunParserConcreteType struct{}

// Real Parser logic
func (k KunParserConcreteType) Parse(in string) (out string) {
	var returnVal = in + "_KUN_parsed"
	return returnVal
}

type DHLParserConcreteType struct{}

// Real Parser logic
func (k DHLParserConcreteType) Parse(in string) (out string) {
	var returnVal = in + "_DHL_parsed"
	return returnVal
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	var handlerConfig = map[string]Service{
		"kun": {
			ClientName: "kun",
			Parser:     &KunParserConcreteType{},
		},
		"DHL": {
			ClientName: "DHL",
			Parser:     &DHLParserConcreteType{},
		},
	}
	var newKunService = handlerConfig["kun"]
	var newDHLService = handlerConfig["DHL"] // DHLParserConcreteType{}

	var d Decider
	d.AddCustomer("kun", newKunService)
	d.AddCustomer("DHL", newDHLService)

	finalHandler := http.HandlerFunc(final)

	mux.Handle("/", d.Decide((finalHandler)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
