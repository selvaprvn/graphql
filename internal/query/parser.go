package query

import (
	"github.com/graphql-go/graphql/language/parser"
	"github.com/selvaprvn/graphql/pkg/logger"
)

// Parse query
func Parse(query string, variables map[string]interface{}) {

	document, err := parser.Parse(parser.ParseParams{Source: query})
	if err != nil {
		logger.Info(err)
	}
	logger.Info(document)
	logger.Info(*document.Loc, string(document.Definitions[0].GetLoc().Source.Body))

}
