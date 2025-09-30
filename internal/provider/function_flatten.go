package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-yamlflattener/internal/flattener"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ function.Function = &flattenFunction{}

// flattenFunction implements the provider function for flattening YAML content
type flattenFunction struct{}

// NewFlattenFunction creates a new instance of the flatten function
func NewFlattenFunction() function.Function {
	return &flattenFunction{}
}

// Metadata returns the function metadata
func (f *flattenFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "flatten"
}

// Definition defines the function signature, parameters, and return type
func (f *flattenFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Flatten nested YAML content into a map with dot notation",
		Description: "Takes YAML content as input and returns a flattened map where nested objects use dot notation (e.g., 'parent.child') and arrays use bracket notation (e.g., 'parent.array[0]').",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "yaml_content",
				Description: "The YAML content to flatten as a string",
			},
		},
		Return: function.MapReturn{
			ElementType: types.StringType,
		},
	}
}

// Run executes the function logic
func (f *flattenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var yamlContent string

	// Get the YAML content parameter
	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &yamlContent))
	if resp.Error != nil {
		return
	}

	// Validate input
	if yamlContent == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("YAML content cannot be empty"))
		return
	}

	// Trim whitespace from YAML content
	yamlContent = strings.TrimSpace(yamlContent)
	if yamlContent == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("YAML content cannot be empty or contain only whitespace"))
		return
	}

	// Check content size for security (10MB limit)
	const maxYAMLSize = 10 * 1024 * 1024
	if len(yamlContent) > maxYAMLSize {
		resp.Error = function.ConcatFuncErrors(resp.Error,
			function.NewFuncError(fmt.Sprintf("YAML content exceeds maximum allowed size of %d bytes", maxYAMLSize)))
		return
	}

	// Sanitize YAML content for security
	yamlContent = strings.ReplaceAll(yamlContent, "\x00", "") // Remove null bytes

	// Create flattener instance with performance and security limits
	flattenerInstance := flattener.Default()
	// Configure flattener with appropriate limits
	flattenerInstance.MaxYAMLSize = maxYAMLSize // 10MB limit
	flattenerInstance.MaxNestingDepth = 100     // Prevent stack overflow
	flattenerInstance.MaxResultSize = 100000    // Limit result size

	flattenedMap, err := flattenerInstance.FlattenYAMLString(yamlContent)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Failed to flatten YAML: "+err.Error()))
		return
	}

	// Convert map[string]string to types.Map
	elements := make(map[string]attr.Value, len(flattenedMap))
	for k, v := range flattenedMap {
		elements[k] = types.StringValue(v)
	}

	resultMap, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Failed to create result map: "+diags[0].Detail()))
		return
	}

	// Set the result
	resp.Result = function.NewResultData(resultMap)
}
