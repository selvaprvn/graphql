package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"github.com/selvaprvn/graphql"
	"github.com/selvaprvn/graphql/handler"
	"github.com/selvaprvn/graphql/pkg/logger"
	"github.com/selvaprvn/graphql/playground"
)

func main() {

	mux := http.NewServeMux()

	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("hello")
	// 	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	// })
	f, err := os.OpenFile("main.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	// msgHand := func(msg string) (string, error) {
	// 	log.Println(msg)

	// 	return `
	// 	{"data": {
	// 		"__schema": {
	// 		  "queryType": {
	// 			"name": "Query"
	// 		  }
	// 		}
	// 	  }
	// 	  }`, nil
	// }
	logg := log.New(f, "", log.LstdFlags)
	logger.SetLogger(logg)
	schema := &graphql.Schema{

		Logger: logg,
		//Handler: msgHand,
	}
	makeCtx := func(ctx context.Context) context.Context {
		return ctx
	}
	logg.Println("start..")
	mux.HandleFunc("/", playground.Handler())
	mux.Handle("/graphql", handler.HTTP(schema, makeCtx))
	mux.Handle("/graphqlws", handler.WebSocket(schema, makeCtx))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
		//Debug:          true,
	})

	log.Fatal(http.ListenAndServe(":8085", c.Handler(mux)))
}
