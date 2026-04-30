package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-yamlflattener/internal/flattener"
)

var _ function.Function = &flattenFunction{}

type flattenFunction struct{}

func NewFlattenFunction() function.Function {
	return &flattenFunction{}
}

func (fn *flattenFunction) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "flatten"
}

func (fn *flattenFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
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

func (fn *flattenFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var yamlContent string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &yamlContent))
	if resp.Error != nil {
		return
	}

	yamlContent = strings.TrimSpace(yamlContent)
	if yamlContent == "" {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("YAML content cannot be empty or contain only whitespace"))
		return
	}

	flattenedMap, err := flattener.Default().FlattenYAMLString(yamlContent)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError(errorTitle(err)+": "+err.Error()))
		return
	}

	resultMap, diags := flattenedToMapValue(flattenedMap)
	if diags.HasError() {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Failed to create result map: "+diags[0].Detail()))
		return
	}

	resp.Result = function.NewResultData(resultMap)
}
