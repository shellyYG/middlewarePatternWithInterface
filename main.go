package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

var key1 = "KEY1"
var key2 = "KEY2"

type Service struct {
	ClientName string
	IParser
	IMapper
}
type IParser interface {
	InnerParse(in string) (out string) // TODO: change to Parser Signature
}
type IMapper interface {
	InnerMap(in string) (out string)
}

func Decider(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Start Deciding...")
		var handlerConfig = map[string]Service{
			"kun": {
				ClientName: "kun",
				IParser:    &KunParserConcreteType{},
				IMapper:    &KunMapperConcreteType{},
			},
		}

		var realTimeClient = "kun"

		var Service = handlerConfig[realTimeClient]
		ctxWithService := context.WithValue(r.Context(), "Service", Service)
		rWithService := r.WithContext(ctxWithService)
		next.ServeHTTP(w, rWithService)
		log.Print("End Deciding.")
	})
}
func Parse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("@Kun Start Parsing...")
		ctx := r.Context()
		svc, ok := ctx.Value("Service").(Service)
		if !ok {
			fmt.Println("oops parsing")
		}
		var in = "一"
		var out = svc.IParser.InnerParse(in)
		ctxWithNewParam1 := context.WithValue(r.Context(), key1, out)
		rWithP1 := r.WithContext(ctxWithNewParam1)
		next.ServeHTTP(w, rWithP1)
		log.Print("@Kun End Parsing.")
	})
}
func Map(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("@Kun Start Mapping...")
		ctx := r.Context()
		svc, ok := ctx.Value("Service").(Service)
		if !ok {
			fmt.Println("oops mapping")
		}
		var in = "一"
		var out = svc.IMapper.InnerMap(in)
		ctxWithNewParam2 := context.WithValue(r.Context(), key2, out)
		rWithP2 := r.WithContext(ctxWithNewParam2)
		next.ServeHTTP(w, rWithP2)
		log.Print("@Kun End Mapping.")
	})
}

type KunParserConcreteType struct{}

func (k *KunParserConcreteType) InnerParse(in string) (out string) {
	var returnVal = in + "_KUN_parsed"
	return returnVal
}

type KunMapperConcreteType struct{}

func (k *KunMapperConcreteType) InnerMap(in string) (out string) {
	var returnVal = in + "_KUN_mapped"
	return returnVal
}

func middlewareThree(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Print("Executing m3")
		p1 := r.Context().Value(key1)
		p2 := r.Context().Value(key2)

		fmt.Println("p1: ", p1, "p2: ", p2)

		if r.URL.Path == "/foo" {
			return
		}
		next.ServeHTTP(w, r)
		log.Print("Executing m3 again")
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Print("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()

	finalHandler := http.HandlerFunc(final)

	mux.Handle("/", Decider(Parse(Map(middlewareThree(finalHandler)))))

	log.Print("Listening on :3000...")
	err := http.ListenAndServe(":3000", mux)
	log.Fatal(err)
}
