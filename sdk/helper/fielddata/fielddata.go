package fielddata

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
)

func Parse(previousValues map[string]interface{}, fieldData *framework.FieldData) (map[string]interface{}, error) {
	if previousValues == nil {
		// This could be provided as nil for convenience if there's no previous value.
		// It will likely be untouched, but populate it just to be defensive
		// against nil pointers.
		previousValues = make(map[string]interface{})
	}
	// Return newValues as a separate map so if the caller needs to compare
	// before/after, they can.
	newValues := make(map[string]interface{})
	var result error
	for schemaName, schema := range fieldData.Schema {
		raw, ok := fieldData.GetOk(schemaName)
		switch {
		case ok:
			// Use the value the user stated.
			newValues[schemaName] = raw
		case previousValues[schemaName] != nil:
			// Retain the previous value.
			newValues[schemaName] = previousValues[schemaName]
		case previousValues[schemaName] == nil:
			// Use the default value held in the schema.
			// (raw holds a nil value so Get is required to pull the default.)
			newValues[schemaName] = fieldData.Get(schemaName)
		}
		// TODO needs testing, not sure if it will work, but this is the general idea
		if newValues[schemaName] == nil && schema.Required {
			result = multierror.Append(result, fmt.Errorf("%q is required but not provided", schemaName))
		}
	}
	return newValues, result
}
