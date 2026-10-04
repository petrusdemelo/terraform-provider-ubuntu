package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = (*ubuntuProvider)(nil)

type ubuntuProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ubuntuProvider{version: version}
	}
}

func (p *ubuntuProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ubuntu"
	resp.Version = p.version
}

func (p *ubuntuProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an Ubuntu host over SSH. Nothing is installed on the host: every operation is a shell command.",
	}
}

func (p *ubuntuProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {
}

func (p *ubuntuProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *ubuntuProvider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}
