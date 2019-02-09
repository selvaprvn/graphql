package schema

// Schema of gql server
type Schema struct {
}

// Register type to schema
func (s *Schema) Register(objectType ...interface{}) *Schema {
	return s
}
