package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-yamlflattener/internal/flattener"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &flattenDataSource{}
var _ datasource.DataSourceWithConfigure = &flattenDataSource{}

type flattenDataSource struct{}

type flattenDataSourceModel struct {
	YAMLContent types.String `tfsdk:"yaml_content"`
	YAMLFile    types.String `tfsdk:"yaml_file"`
	Flattened   types.Map    `tfsdk:"flattened"`
	ID          types.String `tfsdk:"id"`
}

func NewFlattenDataSource() datasource.DataSource {
	return &flattenDataSource{}
}

func (d *flattenDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_flatten"
}

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

func (d *flattenDataSource) Configure(_ context.Context, _ datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
}

func (d *flattenDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data flattenDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.YAMLContent.IsNull() && data.YAMLFile.IsNull() {
		resp.Diagnostics.AddError("Missing Required Input", "Either yaml_content or yaml_file must be provided.")
		return
	}

	if !data.YAMLContent.IsNull() && !data.YAMLFile.IsNull() {
		resp.Diagnostics.AddError("Conflicting Inputs", "Only one of yaml_content or yaml_file should be provided, not both.")
		return
	}

	f := flattener.Default()
	var flattenedMap map[string]string
	var err error

	if !data.YAMLContent.IsNull() {
		flattenedMap, err = f.FlattenYAMLString(data.YAMLContent.ValueString())
	} else {
		flattenedMap, err = f.FlattenYAMLFile(data.YAMLFile.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError(errorTitle(err), err.Error())
		return
	}

	resultMap, diags := flattenedToMapValue(flattenedMap)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	data.Flattened = resultMap
	data.ID = types.StringValue("yaml_flatten")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
