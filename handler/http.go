package handler

import (
	"encoding/json"
	"net/http"

	"github.com/selvaprvn/graphql"
)

// HTTP handler
func HTTP(schema *graphql.Schema, middlewares ...MiddlewareFunc) http.Handler {
	return &httpHandler{
		schema:      schema,
		middlewares: middlewares,
	}
}

type httpHandler struct {
	schema      *graphql.Schema
	middlewares []MiddlewareFunc
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	msg, err := h.schema.Handle(r.Body)
	if err != nil {
		msg = graphql.NewErrorMsg(err)
	}
	json.NewEncoder(w).Encode(msg)
}
