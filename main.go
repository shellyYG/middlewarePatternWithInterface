package main

import (
	"fmt"
	"log"
	"net/http"
)

var key1 = "KEY1"
var key2 = "KEY2"

type Decider struct {
	parser map[string]IParser
}
type IParser interface {
	InnerParse(in string) (out string) // TODO: change to Parser Signature
}

func (d *Decider) AddCustomer(name string, parser IParser) {
	// Initiate map to avoid nil panic
	if d.parser == nil {
		d.parser = make(map[string]IParser)
	}
	d.parser[name] = parser
}

func (d *Decider) Decide(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var realTimeClient = "DHL"
		parser := d.parser[realTimeClient]
		result := parser.InnerParse("Hi") // change to Service.Serve
		fmt.Println("result: ", result)
		next.ServeHTTP(w, r)
		log.Print("End Deciding.")
	})
}

type KunParserConcreteType struct{}

// Real Parser logic
func (k KunParserConcreteType) InnerParse(in string) (out string) {
	var returnVal = in + "_KUN_parsed"
	return returnVal
}

type DHLParserConcreteType struct{}

// Real Parser logic
func (k DHLParserConcreteType) InnerParse(in string) (out string) {
	var returnVal = in + "_DHL_parsed"
	return returnVal
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	var newKunParser = KunParserConcreteType{}
	var newDHLParser = DHLParserConcreteType{}

	var d Decider
	d.AddCustomer("kun", newKunParser)
	d.AddCustomer("DHL", newDHLParser)

	finalHandler := http.HandlerFunc(final)

	mux.Handle("/", d.Decide((finalHandler)))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
