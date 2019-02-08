package graphql

type Schema struct {
	query map[string][]*Schema
}
