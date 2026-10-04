package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrusdemelo/terraform-provider-ubuntu/internal/osrelease"
)

var _ datasource.DataSourceWithConfigure = (*osReleaseDataSource)(nil)

type osReleaseDataSource struct {
	data *providerData
}

type osReleaseModel struct {
	ID              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	PrettyName      types.String `tfsdk:"pretty_name"`
	VersionID       types.String `tfsdk:"version_id"`
	VersionCodename types.String `tfsdk:"version_codename"`
}

func NewOSReleaseDataSource() datasource.DataSource {
	return &osReleaseDataSource{}
}

func (d *osReleaseDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_os_release"
}

func (d *osReleaseDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	computed := func(description string) schema.StringAttribute {
		return schema.StringAttribute{MarkdownDescription: description, Computed: true}
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "The operating system release of the target host, read from `/etc/os-release`.",
		Attributes: map[string]schema.Attribute{
			"id":               computed("The `ID` field, for example `ubuntu`."),
			"name":             computed("The `NAME` field, for example `Ubuntu`."),
			"pretty_name":      computed("The `PRETTY_NAME` field, for example `Ubuntu 24.04.1 LTS`."),
			"version_id":       computed("The `VERSION_ID` field, for example `24.04`."),
			"version_codename": computed("The `VERSION_CODENAME` field, for example `noble`."),
		},
	}
}

func (d *osReleaseDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	data, ok := req.ProviderData.(*providerData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", fmt.Sprintf("expected *providerData, got %T", req.ProviderData))

		return
	}

	d.data = data
}

func (d *osReleaseDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	result, err := d.data.client.Run(ctx, osrelease.Command)
	if err != nil {
		resp.Diagnostics.AddError("Unable to read /etc/os-release", err.Error())

		return
	}

	if result.ExitCode != 0 {
		resp.Diagnostics.AddError("Unable to read /etc/os-release", fmt.Sprintf("exit %d: %s", result.ExitCode, result.Stderr))

		return
	}

	fields := osrelease.Parse(result.Stdout)

	resp.Diagnostics.Append(resp.State.Set(ctx, osReleaseModel{
		ID:              types.StringValue(fields["ID"]),
		Name:            types.StringValue(fields["NAME"]),
		PrettyName:      types.StringValue(fields["PRETTY_NAME"]),
		VersionID:       types.StringValue(fields["VERSION_ID"]),
		VersionCodename: types.StringValue(fields["VERSION_CODENAME"]),
	})...)
}
