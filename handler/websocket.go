package handler

import (
	"context"

	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/selvaprvn/graphql"
)

// MiddlewareFunc middleware
type MiddlewareFunc func(ctx context.Context) context.Context

// WebSocket handle websocket connection
func WebSocket(schema *graphql.Schema, makeCtxFunc MiddlewareFunc) http.Handler {
	upgrader := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		socket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("upgrader.Upgrade: %v", err)
			return
		}
		defer socket.Close()

		makeCtx := func(ctx context.Context) context.Context {
			return makeCtxFunc(ctx)
		}

		log.Println(makeCtx(r.Context()))

		//graphql.ServeJSONSocket(r.Context(), socket, schema, makeCtx, &simpleLogger{})
	})
}
