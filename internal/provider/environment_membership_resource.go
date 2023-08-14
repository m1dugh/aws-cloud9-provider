package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloud9"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/m1dugh/terraform-provider-awscloud9/internal/aws"
)

var _ resource.Resource = &EnvironmentMembershipResource{}

type EnvironmentMembershipResource struct {
	client *aws.AWSCloud9Client
}

type environmentMembershipModel struct {
	EnvironmentId types.String `tfsdk:"environment_id"`
	Permissions   types.String `tfsdk:"permissions"`
	UserARN       types.String `tfsdk:"user_arn"`
}

func NewEnvironmentMembershipResource() resource.Resource {
	return &EnvironmentMembershipResource{}
}

func (rs *EnvironmentMembershipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment_membership"
}

func (rs *EnvironmentMembershipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A membership to a cloud9 environment",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The id of the environment to bound the membership to",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"permissions": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The permissions to give to the role, can be one of `read-write` and `read-only`",
			},
			"user_arn": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The arn of the aws resource that will be given membership to the environment",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (rs *EnvironmentMembershipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aws.AWSCloud9Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure type",
			fmt.Sprintf("Expected *aws.AWSCloud9Client, got %T. Please report this issue to the provider developers.", req.ProviderData))
		return
	}

	rs.client = client
}

func (rs *EnvironmentMembershipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan environmentMembershipModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := plan.EnvironmentId.ValueString()
	_, err := rs.client.Cloud9.CreateEnvironmentMembership(&cloud9.CreateEnvironmentMembershipInput{
		EnvironmentId: &envId,
		UserArn:       plan.UserARN.ValueStringPointer(),
		Permissions:   plan.Permissions.ValueStringPointer(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating environment membership", fmt.Sprintf("An error occured creating membership for environment %s, for user %s: %s", plan.EnvironmentId.String(), plan.UserARN.String(), err.Error()))
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *EnvironmentMembershipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state environmentMembershipModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := state.EnvironmentId.ValueString()
	userArn := state.UserARN.ValueString()

	environments, err := rs.client.GetMemberShips(envId)
	if err != nil {
		resp.Diagnostics.AddError("Error fetching memberships", fmt.Sprintf("Could not retrieve memberships for environment %s", state.EnvironmentId.String()))
	}

	var foundEnv *aws.Cloud9EnvironmentMembership = nil
	for _, env := range environments {
		if env.UserARN == userArn {
			foundEnv = &env
			break
		}
	}

	if foundEnv == nil {
		resp.Diagnostics.AddError("Error fetching memberships", fmt.Sprintf("Could not retrieve membership for environment %s, for user %s", state.EnvironmentId.String(), state.UserARN.String()))
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *EnvironmentMembershipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentMembershipModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := state.EnvironmentId.ValueString()
	userArn := state.UserARN.ValueString()

	_, err := rs.client.Cloud9.DeleteEnvironmentMembership(&cloud9.DeleteEnvironmentMembershipInput{
		EnvironmentId: &envId,
		UserArn:       &userArn,
	})

	if err != nil {
		resp.Diagnostics.AddError("Error deleting membership", fmt.Sprintf("Could not delete membership for environment %s for user %s: %s", envId, state.UserARN.String(), err.Error()))
		return
	}
}

func (rs *EnvironmentMembershipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state environmentMembershipModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	envId := state.EnvironmentId.ValueString()
	userArn := state.UserARN.ValueString()

	_, err := rs.client.Cloud9.UpdateEnvironmentMembership(&cloud9.UpdateEnvironmentMembershipInput{
		EnvironmentId: &envId,
		UserArn:       &userArn,
		Permissions:   state.Permissions.ValueStringPointer(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error updating membership", fmt.Sprintf("Could not update membership for environment %s for user %s: %s", envId, state.UserARN.String(), err.Error()))
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
