package provider

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/m1dugh/terraform-provider-awscloud9/internal/aws"
)

var _ provider.Provider = &AWSCloud9Provider{}

type AWSCloud9Provider struct {
	version string
}

type AWSCloud9ProviderModel struct {
	AccessKeyID     types.String `tfsdk:"aws_access_key_id"`
	SecretAccessKey types.String `tfsdk:"aws_secret_access_key"`
	Region          types.String `tfsdk:"region"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AWSCloud9Provider{
			version: version,
		}
	}
}

func (p *AWSCloud9Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "awscloud9"
	resp.Version = p.version
}

func (p *AWSCloud9Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"aws_access_key_id": schema.StringAttribute{
				MarkdownDescription: "The AWS access key id, if not provided, extracted from `AWS_ACCESS_KEY_ID` env variable.",
				Optional:            true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				MarkdownDescription: "The AWS Secret access key, if not provided, extracted from `AWS_SECRET_ACCESS_KEY` env variable.",
				Optional:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The AWS region to use the provider for, if not provided, extracted from `AWS_REGION` env variable.",
				Optional:            true,
			},
		},
	}
}

func (p *AWSCloud9Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AWSCloud9ProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.AccessKeyID.IsNull() {
		value := os.Getenv("AWS_ACCESS_KEY_ID")
		if len(value) == 0 {
			resp.Diagnostics.AddError("Missing credential", "Missing AWS access key id")
			return
		} else {
			data.AccessKeyID = types.StringValue(value)
		}
	}

	if data.SecretAccessKey.IsNull() {
		value := os.Getenv("AWS_SECRET_ACCESS_KEY")
		if len(value) == 0 {
			resp.Diagnostics.AddError("Missing credential", "Missing AWS secret access key")
			return
		} else {
			data.SecretAccessKey = types.StringValue(value)
		}
	}

	if data.Region.IsNull() {
		value := os.Getenv("AWS_REGION")
		if len(value) == 0 {
			resp.Diagnostics.AddError("Missing aws region", "Missing AWS region configuration")
			return
		} else {
			data.Region = types.StringValue(value)
		}
	}

	client := aws.New(ctx, credentials.NewStaticCredentials(
		data.AccessKeyID.ValueString(), data.SecretAccessKey.ValueString(), ""),
		data.Region.ValueString(),
	)
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AWSCloud9Provider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSSHEnvironmentDataSource,
	}
}

func (p *AWSCloud9Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSSHEnvironmentResource,
		NewEnvironmentMembershipResource,
	}
}
