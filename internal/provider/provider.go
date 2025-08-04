package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	_ "gopkg.in/yaml.v3" // YAML parser dependency
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ provider.Provider = &YAMLFlattenerProvider{}
var _ provider.ProviderWithFunctions = &YAMLFlattenerProvider{}

// YAMLFlattenerProvider defines the provider implementation.
type YAMLFlattenerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// YAMLFlattenerProviderModel describes the provider data model.
type YAMLFlattenerProviderModel struct {
	// Optional configuration fields can be added here if needed in the future
	MaxDepth       types.Int64 `tfsdk:"max_depth"`
	EscapeNewlines types.Bool  `tfsdk:"escape_newlines"`
}

// Metadata returns the provider metadata including type name and version.
func (p *YAMLFlattenerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "yamlflattener"
	resp.Version = p.version

	// Description is set in the schema, not here
}

// Schema defines the provider schema and configuration options.
func (p *YAMLFlattenerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The YAML Flattener provider allows you to flatten nested YAML structures into flat key-value maps with dot notation for nested objects and bracket notation for arrays.",
		Attributes: map[string]schema.Attribute{
			"max_depth": schema.Int64Attribute{
				Description: "Maximum recursion depth for flattening (default: 100). Set to prevent stack overflow with deeply nested structures.",
				Optional:    true,
			},
			"escape_newlines": schema.BoolAttribute{
				Description: "When true, newlines in multi-line values are escaped as \\n for compatibility with tools that parse values as key-value pairs (default: false).",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{},
	}
}

// Configure configures the provider with user-provided configuration.
func (p *YAMLFlattenerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data YAMLFlattenerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Pass configuration to data sources and functions
	resp.DataSourceData = &data
	resp.ResourceData = &data
}

// Resources returns the list of resources supported by this provider.
func (p *YAMLFlattenerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// No resources in this provider
	}
}

// DataSources returns the list of data sources supported by this provider.
func (p *YAMLFlattenerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewFlattenDataSource,
	}
}

// Functions returns the list of functions supported by this provider.
func (p *YAMLFlattenerProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewFlattenFunction,
	}
}

// New creates a new instance of the YAML Flattener provider
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &YAMLFlattenerProvider{
			version: version,
		}
	}
}
