package main

import (
	"context"
	"fmt"
	"html"
	"log"
	"net/http"

	"github.com/selvaprvn/graphql"
	"github.com/selvaprvn/graphql/handler"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	schema := &graphql.Schema{}
	http.Handle("/graphqlws", handler.WebSocket(schema, func(ctx context.Context) {
		return ctx
	}))
	log.Fatal(http.ListenAndServe(":8085", nil))
}
