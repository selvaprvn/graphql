package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gorilla/websocket"
	"github.com/selvaprvn/graphql"
)

type messageType string

const (
	msgEcho        messageType = "echo"
	msgError       messageType = "error"
	msgSubscribe   messageType = "subscribe"
	msgUnSubscribe messageType = "unsubscribe"
	msgMutate      messageType = "mutate"
)

type jsonSocket interface {
	ReadJSON(value interface{}) error
	WriteJSON(value interface{}) error
	Close() error
}

// MiddlewareFunc middleware
type MiddlewareFunc func(ctx context.Context) context.Context

// WebSocket handle websocket connection
func WebSocket(schema *graphql.Schema, middlewares ...MiddlewareFunc) http.Handler {
	return &wsHandler{
		httpHandler: httpHandler{
			schema:      schema,
			middlewares: middlewares,
		},
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

type wsHandler struct {
	httpHandler
	upgrader *websocket.Upgrader
}

func (h *wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("upgrader.Upgrade: %v", err), http.StatusInternalServerError)
		return
	}

	wsSoc := &wsSocket{
		socket:      socket,
		httpHandler: h.httpHandler,
	}

	defer wsSoc.close()

	wsSoc.run()

}

type wsSocket struct {
	httpHandler
	socket *websocket.Conn
	closed bool
}

type inMessage struct {
	ID         string                 `json:"id"`
	Type       messageType            `json:"type"`
	Message    json.RawMessage        `json:"message"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

type outMessage struct {
	ID       string                 `json:"id,omitempty"`
	Type     messageType            `json:"type"`
	Message  interface{}            `json:"message,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// newErrorMsg new error out message
func newErrorMsg(id string, err error) *outMessage {
	return &outMessage{
		ID:       id,
		Type:     msgError,
		Message:  err.Error(),
		Metadata: nil,
	}
}

func (ws *wsSocket) read() (*inMessage, bool) {
	var inMsg inMessage
	if err := ws.socket.ReadJSON(&inMsg); err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Printf("error: %v", err)
		} else {
			ws.close()
		}
		return nil, true
	}
	return &inMsg, false
}

func (ws *wsSocket) write(msg *outMessage) bool {
	if err := ws.socket.WriteJSON(msg); err != nil {
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			fmt.Printf("error: %v", err)
		} else {
			ws.close()
		}
		return true
	}
	return false
}

func (ws *wsSocket) close() {
	ws.closed = true
	ws.socket.Close()
}

func (ws *wsSocket) run() {
	for {
		msg, closed := ws.read()
		if closed || ws.closed {
			break
		}

		go func() {

			gmsg, err := ws.schema.Handle(bytes.NewReader(msg.Message))
			if err != nil {
				ws.write(newErrorMsg(msg.ID, err))
				return
			}
			resp := &outMessage{
				ID:      msg.ID,
				Message: gmsg,
			}
			ws.write(resp)
		}()

	}
}
