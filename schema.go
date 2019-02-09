package graphql

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"

	"github.com/selvaprvn/graphql/internal/query"
)

type operationType string

const (
	introspectionQuery operationType = "IntrospectionQuery"
)

// MessageHandler middleware
type MessageHandler func(msg string) (string, error)

// Schema gql schema
type Schema struct {
	Handler MessageHandler
	Logger  *log.Logger
}

// GqlRequest input message
type GqlRequest struct {
	OperationName operationType          `json:"operationName"`
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
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
		s.Logger.Println("reqs:", string(byt))
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
			s.Logger.Println("parsing")
			var inMsg = new(GqlRequest)
			if err := json.Unmarshal(byt, inMsg); err != nil {
				gResp = NewErrorMsg(err)
			}
			gResp = s.handle(inMsg)
		}
	}
	return gResp, nil
}

func (s *Schema) handle(req *GqlRequest) *GqlResponse {
	//log.Println(req)
	gresp := &GqlResponse{}
	s.Logger.Println("req", req.OperationName, req)
	query.Parse(req.Query, req.Variables)
	if req.OperationName == introspectionQuery {
		gresp.Data = map[string]interface{}{
			"__schema": map[string]interface{}{
				"directives": []map[string]interface{}{},
				"mutationType": map[string]interface{}{
					"name": "Mutation",
				},
				"queryType": map[string]interface{}{
					"name": "Query",
				},
				"subscriptionType": map[string]interface{}{
					"name": "Subscription",
				},
				"types": []map[string]interface{}{
					{"name": "Query",
						"kind": "OBJECT",
						"fields": []map[string]interface{}{{
							"name": "id",
							"type": map[string]interface{}{
								"name": "null",
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"name": "ID",
									"kind": "SCALAR",
								},
							},
						}},
						"interfaces": []map[string]interface{}{},

						"possibleTypes": "null",
					},
				},
			},
		}
	}
	s.Logger.Println(gresp)
	return gresp
}
