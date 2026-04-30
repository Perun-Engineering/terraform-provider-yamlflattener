package provider

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-yamlflattener/internal/flattener"
)

var errorTitles = map[flattener.ErrorType]string{
	flattener.ErrTypeValidation: "Invalid Input",
	flattener.ErrTypeParsing:    "Invalid YAML Syntax",
	flattener.ErrTypeDepthLimit: "Nesting Depth Exceeded",
	flattener.ErrTypeSizeLimit:  "Size Limit Exceeded",
	flattener.ErrTypeTimeout:    "Operation Timed Out",
	flattener.ErrTypeFileAccess: "File Access Error",
	flattener.ErrTypeSecurity:   "Security Error",
}

// errorTitle returns a human-readable title for a flattener error, or "Flatten Error" for unknown errors.
func errorTitle(err error) string {
	var fe *flattener.Error
	if errors.As(err, &fe) {
		if title := errorTitles[fe.Type]; title != "" {
			return title
		}
	}
	return "Flatten Error"
}

// flattenedToMapValue converts a map[string]string to a Terraform types.Map.
func flattenedToMapValue(m map[string]string) (types.Map, diag.Diagnostics) {
	elements := make(map[string]attr.Value, len(m))
	for k, v := range m {
		elements[k] = types.StringValue(v)
	}
	return types.MapValue(types.StringType, elements)
}
