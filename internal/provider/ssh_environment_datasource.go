package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/m1dugh/terraform-provider-awscloud9/internal/aws"
)

var _ datasource.DataSource = &SSHEnvironmentDataSource{}

func NewSSHEnvironmentDataSource() datasource.DataSource {
	return &SSHEnvironmentDataSource{}
}

type SSHEnvironmentDataSource struct {
	client *aws.AWSCloud9Client
}

type membershipModel struct {
	Permissions types.String `tfsdk:"permissions"`
	UserARN     types.String `tfsdk:"user_arn"`
}

type SSHEnvironmentDataSourceModel struct {
	Arn             types.String `tfsdk:"arn"`
	EnvironmentId   types.String `tfsdk:"environment_id"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	LoginName       types.String `tfsdk:"login_name"`
	Hostname        types.String `tfsdk:"hostname"`
	Port            types.Int64  `tfsdk:"port"`
	EnvironmentPath types.String `tfsdk:"environment_path"`
	NodePath        types.String `tfsdk:"node_path"`
	BastionURL      types.String `tfsdk:"bastion_url"`
	Tags            types.Map    `tfsdk:"tags"`
}

func (ds *SSHEnvironmentDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_environment"
}

func (ds *SSHEnvironmentDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieves a cloud 9 SSH environment",
		Attributes: map[string]schema.Attribute{
			"arn": schema.StringAttribute{
				MarkdownDescription: "The ARN of the environment",
				Optional:            false,
				Required:            false,
				Computed:            true,
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The id of the cloud 9 environment",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the environment",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the environment",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"login_name": schema.StringAttribute{
				MarkdownDescription: "The login name of the user bound to the environment",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname to connect to",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "The port to connect to",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"environment_path": schema.StringAttribute{
				MarkdownDescription: "The path where the cloud 9 environment shoud open a shell into.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"node_path": schema.StringAttribute{
				MarkdownDescription: "The path where node is set on the remote host.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"bastion_url": schema.StringAttribute{
				MarkdownDescription: "The url to connect to the bastion host.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "The tags of the environment.",
				ElementType:         types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
}

func (ds *SSHEnvironmentDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aws.AWSCloud9Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aws.AWSCloud9Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	ds.client = client
}

func (ds *SSHEnvironmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data SSHEnvironmentDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	environmentId := data.EnvironmentId.ValueString()
	environments, err := ds.client.GetSSHEnvironments(environmentId)
	if err != nil {
		resp.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to read environment %s, got error: %s", environmentId, err))
		return
	} else if len(environments) == 0 {
		resp.Diagnostics.AddError("Environment not found", fmt.Sprintf("Unable to read environment %s", environmentId))
	}

	environment := environments[0]
	diags := convertModelToPlan(&data, &environment)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
