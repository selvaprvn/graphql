package query

import "testing"

func TestParse(t *testing.T) {
	type args struct {
		query     string
		variables map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "IntrospectionQuery",
			args: args{
				query:     IntrospectionQuery,
				variables: map[string]interface{}{},
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Parse(tt.args.query, tt.args.variables)
		})
	}
}

const IntrospectionQuery = `
query IntrospectionQuery {
	__schema {
	  queryType {
		name
	  }
	  mutationType {
		name
	  }
	  subscriptionType {
		name
	  }
	  types {
		...FullType
	  }
	  directives {
		name
		description
		locations
		args {
		  ...InputValue
		}
	  }
	}
  }
  
  fragment FullType on __Type {
	kind
	name
	description
	fields(includeDeprecated: true) {
	  name
	  description
	  args {
		...InputValue
	  }
	  type {
		...TypeRef
	  }
	  isDeprecated
	  deprecationReason
	}
	inputFields {
	  ...InputValue
	}
	interfaces {
	  ...TypeRef
	}
	enumValues(includeDeprecated: true) {
	  name
	  description
	  isDeprecated
	  deprecationReason
	}
	possibleTypes {
	  ...TypeRef
	}
  }
  
  fragment InputValue on __InputValue {
	name
	description
	type {
	  ...TypeRef
	}
	defaultValue
  }
  
  fragment TypeRef on __Type {
	kind
	name
	ofType {
	  kind
	  name
	  ofType {
		kind
		name
		ofType {
		  kind
		  name
		  ofType {
			kind
			name
			ofType {
			  kind
			  name
			  ofType {
				kind
				name
				ofType {
				  kind
				  name
				}
			  }
			}
		  }
		}
	  }
	}
  }
`
