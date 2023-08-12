package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/cloud9"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/m1dugh/terraform-provider-awscloud9/internal/aws"
)

var _ resource.Resource = &SSHEnvironmentResource{}

type SSHEnvironmentResource struct {
	client *aws.AWSCloud9Client
}

func NewSSHEnvironmentResource() resource.Resource {
	return &SSHEnvironmentResource{}
}

type SSHEnvironmentResourceModel = SSHEnvironmentDataSourceModel

func (rs *SSHEnvironmentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_environment"
}

func (rs *SSHEnvironmentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Creates a cloud 9 SSH environment",
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required:            false,
				Optional:            false,
				Computed:            true,
				MarkdownDescription: "The id of the environment",
                PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the environment",
			},
            "description": schema.StringAttribute{
                Required: false,
                Optional: true,
                MarkdownDescription: "The description of the environment",
            },
			"login_name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The login name of the user to use the environment",
			},
			"hostname": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The hostname of the remote machine",
			},
			"port": schema.Int64Attribute{
				Default:  int64default.StaticInt64(22),
				Computed: true,
			},
			"environment_path": schema.StringAttribute{
				Required:            false,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The path for the environment",
                PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"node_path": schema.StringAttribute{
				Required:            false,
				Optional:            true,
                Computed: true,
				MarkdownDescription: "The path to node.js on the remote host",
                PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"bastion_url": schema.StringAttribute{
				Required:            false,
				Optional:            true,
				MarkdownDescription: "The ssh url to a bastion host",
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "A list of tags to attach",
				ElementType:         types.StringType,
				Required:            false,
				Optional:            true,
			},
            "memberships": schema.ListNestedAttribute {
                MarkdownDescription: "The memberships to bind to the environment",
                Required: false,
                Optional: true,
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "user_arn": schema.StringAttribute{
                            Required: true,
                        },
                        "permissions": schema.ListAttribute{
                            ElementType: types.StringType,
                            Required: false,
                            Optional: true,
                        },
                    },
                },  
            },
		},
	}
}

func (rs *SSHEnvironmentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	rs.client = client
}

func (rs *SSHEnvironmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan SSHEnvironmentResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var request aws.CreateEnvironmentSSHRequest

	request.Name = plan.Name.ValueString()
	request.LoginName = plan.LoginName.ValueString()
	request.Hostname = plan.Hostname.ValueString()
	request.Port = int16(plan.Port.ValueInt64())

	if !plan.BastionURL.IsNull() {
		request.BastionHost = plan.BastionURL.ValueString()
	}

	if !plan.NodePath.IsNull() {
		request.NodePath = plan.NodePath.ValueString()
	}

	if !plan.EnvironmentPath.IsNull() {
		request.EnvironmentPath = plan.EnvironmentId.ValueString()
	} else {
        plan.EnvironmentPath = types.StringValue("")
    }

	environment, err := rs.client.CreateEnvironmentSSH(&request)
	if err != nil {
		resp.Diagnostics.AddError("Client error", fmt.Sprintf("Unable to create environment %s, got error: %s", request.Name, err))
        return
	}

	plan.EnvironmentId = types.StringValue(environment.EnvironmentId)
    plan.Memberships = basetypes.NewListNull(types.ObjectType{
        AttrTypes: map[string]attr.Type{
            "user_arn": types.StringType,
            "permissions": types.ListType{
                ElemType: types.StringType,
            },
        },
    })

    readResult, err := rs.client.DescribeSSHRemote(environment.EnvironmentId)
    if err != nil {
        resp.Diagnostics.AddError("Client error", fmt.Sprintf("Could not read environment %s: %s", environment.EnvironmentId, err.Error()))
    }
    plan.NodePath = types.StringValue(readResult.Results.NodePath)
    plan.EnvironmentPath = types.StringValue(readResult.Results.EnvironmentPath)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (rs *SSHEnvironmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state SSHEnvironmentResourceModel

    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    envId := state.EnvironmentId.ValueString()
    environments, err := rs.client.GetSSHEnvironments([]string{envId})
    if err != nil {
        resp.Diagnostics.AddError("Error fetching env", fmt.Sprintf("Could not fetch env %s: %s", envId, err.Error()))
        return
    }

    if len(environments) == 0 {
        resp.Diagnostics.AddError("Environment not found", fmt.Sprintf("Environment not found: %s", envId))
        return
    }

    environment := environments[0]
    state.EnvironmentId = basetypes.NewStringValue(environment.EnvironmentId)
    state.BastionURL = basetypes.NewStringValue(environment.BastionHost)
    state.Name = basetypes.NewStringValue(environment.Name)
    state.LoginName = basetypes.NewStringValue(environment.LoginName)
    state.Description = basetypes.NewStringValue(environment.Description)
    state.Port = basetypes.NewInt64Value(int64(environment.Port))
    state.Hostname = basetypes.NewStringValue(environment.Hostname)
    state.EnvironmentPath = basetypes.NewStringValue(environment.EnvironmentPath)
    state.NodePath = basetypes.NewStringValue(environment.NodePath)
    typedTags := make(map[string]attr.Value)
    for key, val := range environment.Tags {
        typedTags[key] = basetypes.NewStringValue(val)
    }
    state.Tags, diags = basetypes.NewMapValue(types.StringType, typedTags)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    state.Memberships = basetypes.NewListNull(types.ObjectType{
        AttrTypes: map[string]attr.Type{
            "user_arn": types.StringType,
            "permissions": types.ListType{
                ElemType: types.StringType,
            },
        },
    })

    diags = resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}

func (rs *SSHEnvironmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state SSHEnvironmentResourceModel
    diags := req.State.Get(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    envId := state.EnvironmentId.ValueString()
    _, err := rs.client.Cloud9.DeleteEnvironment(&cloud9.DeleteEnvironmentInput{
        EnvironmentId: &envId,
    })
    if err != nil {
        resp.Diagnostics.AddError("Error deleting env", fmt.Sprintf("Could not delete environment %s: %s", envId, err.Error()))
        return
    }
}

func (rs *SSHEnvironmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan SSHEnvironmentResourceModel
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    var updateRequest aws.UpdateSSHRemoteRequest
    updateRequest.Hostname = plan.Hostname.ValueString()
    envId := plan.EnvironmentId.ValueString()
    updateRequest.LoginName = plan.LoginName.ValueString()
    updateRequest.EnvironmentId = envId
    if !plan.Port.IsNull() {
        updateRequest.Port = int16(plan.Port.ValueInt64())
    }
    if !plan.EnvironmentPath.IsNull() {
        updateRequest.EnvironmentPath = plan.EnvironmentPath.ValueString()
    }
    if !plan.NodePath.IsNull() {
        updateRequest.NodePath = plan.NodePath.ValueString()
    }

    err := rs.client.UpdateSSHRemote(&updateRequest)
    if err != nil {
        resp.Diagnostics.AddError("Error updating environment", fmt.Sprintf("Error updating environment %s: %s", envId, err.Error()))
        return
    }

    plan.Memberships = basetypes.NewListNull(types.ObjectType{
        AttrTypes: map[string]attr.Type{
            "user_arn": types.StringType,
            "permissions": types.ListType{
                ElemType: types.StringType,
            },
        },
    })

    diags = resp.State.Set(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
}
