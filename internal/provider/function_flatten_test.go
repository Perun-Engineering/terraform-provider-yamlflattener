package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestFlattenFunction_Metadata(t *testing.T) {
	f := NewFlattenFunction(nil)

	resp := &function.MetadataResponse{}
	f.Metadata(context.Background(), function.MetadataRequest{}, resp)

	if resp.Name != "flatten" {
		t.Errorf("Expected function name 'flatten', got %s", resp.Name)
	}
}

func TestFlattenFunction_Definition(t *testing.T) {
	f := NewFlattenFunction(nil)

	resp := &function.DefinitionResponse{}
	f.Definition(context.Background(), function.DefinitionRequest{}, resp)

	if len(resp.Definition.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(resp.Definition.Parameters))
	}

	if resp.Definition.Parameters[0].GetName() != "yaml_content" {
		t.Errorf("Expected parameter name 'yaml_content', got %s", resp.Definition.Parameters[0].GetName())
	}
}

func TestFlattenFunction_Run_EmptyContent(t *testing.T) {
	f := NewFlattenFunction(nil)

	resp := &function.RunResponse{}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.StringValue("")}),
	}, resp)

	if resp.Error == nil {
		t.Error("Expected error for empty YAML content, got nil")
	}
}

func TestFlattenFunction_Run_WhitespaceOnly(t *testing.T) {
	f := NewFlattenFunction(nil)

	resp := &function.RunResponse{}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.StringValue("   \n   \t   ")}),
	}, resp)

	if resp.Error == nil {
		t.Error("Expected error for whitespace-only YAML content, got nil")
	}
}

func TestFlattenFunction_Run_InvalidYAML(t *testing.T) {
	f := NewFlattenFunction(nil)

	resp := &function.RunResponse{}
	f.Run(context.Background(), function.RunRequest{
		Arguments: function.NewArgumentsData([]attr.Value{types.StringValue("key1: value1\nkey2: [\n  invalid yaml\n")}),
	}, resp)

	if resp.Error == nil {
		t.Error("Expected error for invalid YAML content, got nil")
	}
}
