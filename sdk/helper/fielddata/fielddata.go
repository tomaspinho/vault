package fielddata

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
)

func CreateUpdate(previousValues map[string]interface{}, fieldData *framework.FieldData) (map[string]interface{}, error) {
	if previousValues == nil {
		// This can be provided as nil for convenience if there's no previous value.
		previousValues = make(map[string]interface{})
	}
	// Return newValues as a separate map so if the caller needs to compare before/after, they can.
	newValues := make(map[string]interface{})
	var errs error
	for schemaName, schema := range fieldData.Schema {
		raw, ok := fieldData.GetOk(schemaName)
		switch {
		case ok:
			// Use the value the user stated.
			newValues[schemaName] = raw
		case previousValues[schemaName] != nil:
			// Retain the previous value.
			newValues[schemaName] = previousValues[schemaName]
		case schema.Default != nil:
			// Use the default value held in the schema.
			newValues[schemaName] = schema.Default
		case schema.Required:
			errs = multierror.Append(errs, fmt.Errorf("%q is required but not provided", schemaName))
		}
	}
	return newValues, errs
}
