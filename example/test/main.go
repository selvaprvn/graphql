package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"

	//"github.com/rs/cors"
	"github.com/selvaprvn/graphql"
	"github.com/selvaprvn/graphql/handler"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("hello")
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	msgHand := func(msg string) (string, error) {
		log.Println(msg)
		return "", nil
	}
	schema := &graphql.Schema{
		Handler: msgHand,
	}
	makeCtx := func(ctx context.Context) context.Context {
		return ctx
	}
	mux.Handle("/graphql", handler.HTTP(schema, makeCtx))
	mux.Handle("/graphqlws", handler.WebSocket(schema, makeCtx))

	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedHeaders: []string{"*"},
	// 	Debug:          true,
	// })

	log.Fatal(http.ListenAndServe(":8085", mux))
}
