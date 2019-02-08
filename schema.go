package graphql

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// MessageHandler middleware
type MessageHandler func(msg string) (string, error)

// Schema gql schema
type Schema struct {
	Handler MessageHandler
}

// GqlRequest input message
type GqlRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GqlResponse output message
type GqlResponse struct {
	Data   interface{} `json:"data"`
	Errors []string    `json:"errors"`
}

// NewErrorMsg new gql err
func NewErrorMsg(errs ...error) *GqlResponse {
	resp := &GqlResponse{}
	for _, e := range errs {
		resp.Errors = append(resp.Errors, e.Error())
	}
	return resp
}

// Handle messages
func (s *Schema) Handle(reader io.Reader) (*GqlResponse, error) {
	var gResp *GqlResponse
	byt, rerr := ioutil.ReadAll(reader)
	if rerr == nil {
		if s.Handler != nil {
			sresp, er := s.Handler(string(byt))
			if er != nil {
				gResp = NewErrorMsg(er)
			} else {
				if err := json.Unmarshal([]byte(sresp), gResp); err != nil {
					gResp = NewErrorMsg(err)
				}
			}
		} else {
			var inMsg *GqlRequest
			if err := json.Unmarshal([]byte(byt), inMsg); err != nil {
				gResp = NewErrorMsg(err)
			}
			gResp = s.handle(inMsg)
		}
	}
	return gResp, nil
}

func (s *Schema) handle(req *GqlRequest) *GqlResponse {
	return &GqlResponse{}
}
