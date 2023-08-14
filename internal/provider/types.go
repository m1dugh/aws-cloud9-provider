package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

type SSHEnvironmentModel struct {
	Arn             types.String `tfsdk:"arn"`
	ID              types.String `tfsdk:"id"`
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
