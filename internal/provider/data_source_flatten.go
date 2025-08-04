package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-yamlflattener/internal/flattener"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &flattenDataSource{}
var _ datasource.DataSourceWithConfigure = &flattenDataSource{}

// flattenDataSource defines the data source implementation.
type flattenDataSource struct {
	// Provider configuration
	providerData *YAMLFlattenerProviderModel
}

// flattenDataSourceModel describes the data source data model.
type flattenDataSourceModel struct {
	YAMLContent types.String `tfsdk:"yaml_content"`
	YAMLFile    types.String `tfsdk:"yaml_file"`
	Flattened   types.Map    `tfsdk:"flattened"`
	ID          types.String `tfsdk:"id"`
}

// NewFlattenDataSource creates a new instance of the flatten data source
func NewFlattenDataSource() datasource.DataSource {
	return &flattenDataSource{}
}

// Metadata returns the data source type name.
func (d *flattenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flatten"
}

// Schema defines the schema for the data source.
func (d *flattenDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Flattens nested YAML structures into a map with dot notation for nested objects and bracket notation for arrays.",
		Attributes: map[string]schema.Attribute{
			"yaml_content": schema.StringAttribute{
				Description: "YAML content to flatten as a string. Either yaml_content or yaml_file must be provided.",
				Optional:    true,
			},
			"yaml_file": schema.StringAttribute{
				Description: "Path to a YAML file to flatten. Either yaml_content or yaml_file must be provided.",
				Optional:    true,
			},
			"flattened": schema.MapAttribute{
				Description: "The resulting flattened map where nested objects use dot notation and arrays use bracket notation.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"id": schema.StringAttribute{
				Description: "Identifier for this data source instance.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *flattenDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*YAMLFlattenerProviderModel)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *YAMLFlattenerProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.providerData = providerData
}

// Read refreshes the Terraform state with the latest data.
func (d *flattenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data flattenDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate input: either yaml_content or yaml_file must be provided, but not both
	if data.YAMLContent.IsNull() && data.YAMLFile.IsNull() {
		resp.Diagnostics.AddError(
			"Missing Required Input",
			"Either yaml_content or yaml_file must be provided.",
		)
		return
	}

	if !data.YAMLContent.IsNull() && !data.YAMLFile.IsNull() {
		resp.Diagnostics.AddError(
			"Conflicting Inputs",
			"Only one of yaml_content or yaml_file should be provided, not both.",
		)
		return
	}

	// Create flattener instance with performance and security limits
	escapeNewlines := false
	if d.providerData != nil && !d.providerData.EscapeNewlines.IsNull() {
		escapeNewlines = d.providerData.EscapeNewlines.ValueBool()
	}

	flattenerInstance := flattener.NewFlattenerWithOptions(escapeNewlines)
	// Configure flattener with appropriate limits
	flattenerInstance.MaxYAMLSize = 10 * 1024 * 1024 // 10MB limit
	flattenerInstance.MaxNestingDepth = 100          // Prevent stack overflow
	flattenerInstance.MaxResultSize = 100000         // Limit result size

	// Apply provider configuration
	if d.providerData != nil && !d.providerData.MaxDepth.IsNull() {
		flattenerInstance.MaxNestingDepth = int(d.providerData.MaxDepth.ValueInt64())
	}

	var flattenedMap map[string]string
	var err error

	// Process based on input type with additional security measures
	if !data.YAMLContent.IsNull() {
		yamlContent := data.YAMLContent.ValueString()

		// Check content size
		if len(yamlContent) > flattenerInstance.MaxYAMLSize {
			resp.Diagnostics.AddError(
				"YAML Content Too Large",
				fmt.Sprintf("YAML content exceeds maximum allowed size of %d bytes", flattenerInstance.MaxYAMLSize),
			)
			return
		}

		flattenedMap, err = flattenerInstance.FlattenYAMLString(yamlContent)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Flatten YAML Content",
				fmt.Sprintf("Error flattening YAML content: %s", err),
			)
			return
		}
	} else {
		yamlFile := data.YAMLFile.ValueString()

		// Validate file path for security
		if strings.Contains(filepath.Clean(yamlFile), "..") {
			resp.Diagnostics.AddError(
				"Invalid File Path",
				"File path contains directory traversal patterns which are not allowed for security reasons.",
			)
			return
		}

		flattenedMap, err = flattenerInstance.FlattenYAMLFile(yamlFile)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Flatten YAML File",
				fmt.Sprintf("Error flattening YAML file %s: %s", yamlFile, err),
			)
			return
		}
	}

	// Convert map[string]string to types.Map
	elements := make(map[string]attr.Value, len(flattenedMap))
	for k, v := range flattenedMap {
		elements[k] = types.StringValue(v)
	}

	resultMap, diags := types.MapValue(types.StringType, elements)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set the flattened result
	data.Flattened = resultMap

	// Set a unique ID for the data source
	data.ID = types.StringValue("yaml_flatten")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
