package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/petrusdemelo/terraform-provider-ubuntu/internal/ssh"
)

const (
	defaultPort = int64(22)
	defaultSudo = true
)

var _ provider.Provider = (*ubuntuProvider)(nil)

type ubuntuProvider struct {
	version string
}

type providerModel struct {
	SSH           *sshModel    `tfsdk:"ssh"`
	DefaultTarget *targetModel `tfsdk:"default_target"`
	Sudo          types.Bool   `tfsdk:"sudo"`
}

type sshModel struct {
	User                  types.String `tfsdk:"user"`
	PrivateKey            types.String `tfsdk:"private_key"`
	Password              types.String `tfsdk:"password"`
	HostKey               types.String `tfsdk:"host_key"`
	InsecureIgnoreHostKey types.Bool   `tfsdk:"insecure_ignore_host_key"`
}

type targetModel struct {
	Target types.String `tfsdk:"target"`
	Port   types.Int64  `tfsdk:"port"`
}

type commandRunner interface {
	Run(ctx context.Context, cmd string) (ssh.Result, error)
}

type providerData struct {
	client commandRunner
	sudo   bool
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
		Attributes: map[string]schema.Attribute{
			"sudo": schema.BoolAttribute{
				MarkdownDescription: "Run privileged commands through `sudo`. Defaults to `true`.",
				Optional:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"ssh": schema.SingleNestedBlock{
				MarkdownDescription: "How to authenticate the SSH connection.",
				Validators:          []validator.Object{objectvalidator.IsRequired()},
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						MarkdownDescription: "SSH user.",
						Optional:            true,
					},
					"private_key": schema.StringAttribute{
						MarkdownDescription: "PEM-encoded private key. Takes precedence over `password`.",
						Optional:            true,
						Sensitive:           true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Password, used when `private_key` is unset.",
						Optional:            true,
						Sensitive:           true,
					},
					"host_key": schema.StringAttribute{
						MarkdownDescription: "Expected host public key, in `authorized_keys` format. Required unless `insecure_ignore_host_key` is `true`.",
						Optional:            true,
					},
					"insecure_ignore_host_key": schema.BoolAttribute{
						MarkdownDescription: "Skip host key verification. Defaults to `false`.",
						Optional:            true,
					},
				},
			},
			"default_target": schema.SingleNestedBlock{
				MarkdownDescription: "The host to manage.",
				Validators:          []validator.Object{objectvalidator.IsRequired()},
				Attributes: map[string]schema.Attribute{
					"target": schema.StringAttribute{
						MarkdownDescription: "Hostname or IP address.",
						Optional:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "SSH port. Defaults to `22`.",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (p *ubuntuProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.SSH == nil || config.DefaultTarget == nil {
		resp.Diagnostics.AddError("Missing provider configuration", "Both the ssh and default_target blocks are required.")

		return
	}

	for _, value := range []types.String{config.DefaultTarget.Target, config.SSH.User, config.SSH.PrivateKey, config.SSH.Password, config.SSH.HostKey} {
		if value.IsUnknown() {
			resp.Diagnostics.AddError(
				"Unknown provider configuration",
				"The connection settings must be known at plan time. Apply the resource they come from first, or pass them in as variables.",
			)

			return
		}
	}

	if config.DefaultTarget.Target.ValueString() == "" || config.SSH.User.ValueString() == "" {
		resp.Diagnostics.AddError("Missing provider configuration", "default_target.target and ssh.user must be set.")

		return
	}

	port := defaultPort
	if !config.DefaultTarget.Port.IsNull() {
		port = config.DefaultTarget.Port.ValueInt64()
	}

	sudo := defaultSudo
	if !config.Sudo.IsNull() {
		sudo = config.Sudo.ValueBool()
	}

	client, err := ssh.Dial(ctx, ssh.Config{
		Host:                  config.DefaultTarget.Target.ValueString(),
		Port:                  port,
		User:                  config.SSH.User.ValueString(),
		PrivateKey:            config.SSH.PrivateKey.ValueString(),
		Password:              config.SSH.Password.ValueString(),
		HostKey:               config.SSH.HostKey.ValueString(),
		InsecureIgnoreHostKey: config.SSH.InsecureIgnoreHostKey.ValueBool(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Unable to connect to the Ubuntu host", err.Error())

		return
	}

	data := &providerData{client: client, sudo: sudo}
	resp.ResourceData = data
	resp.DataSourceData = data
}

func (p *ubuntuProvider) Resources(context.Context) []func() resource.Resource {
	return nil
}

func (p *ubuntuProvider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOSReleaseDataSource,
	}
}
