package provider

import (
	"context"
	"os"

	"github.com/Devolutions/go-dvls"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure DvlsProvider satisfies various provider interfaces.
var _ provider.Provider = &DvlsProvider{}
var _ provider.ProviderWithEphemeralResources = &DvlsProvider{}

// DvlsProvider defines the provider implementation.
type DvlsProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// DvlsProviderModel describes the provider data model.
type DvlsProviderModel struct {
	BaseUri   types.String `tfsdk:"base_uri"`
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
	ApiKey    types.String `tfsdk:"api_key"`
}

func (p *DvlsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "dvls"
	resp.Version = p.version
}

func (p *DvlsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The provider authenticates either with an application (`app_id` + `app_secret`) or with an API key (`api_key`). " +
			"Credentials fall back to the environment variables `DVLS_APP_ID`, `DVLS_APP_SECRET` and `DVLS_API_KEY`. " +
			"An `api_key` set in the configuration is always used; otherwise application credentials take precedence over `DVLS_API_KEY`.",
		Attributes: map[string]schema.Attribute{
			"base_uri": schema.StringAttribute{
				Description: "DVLS base URI",
				Required:    true,
			},
			"app_id": schema.StringAttribute{
				MarkdownDescription: "DVLS App ID `$DVLS_APP_ID`",
				Optional:            true,
			},
			"app_secret": schema.StringAttribute{
				MarkdownDescription: "DVLS App Secret `$DVLS_APP_SECRET`",
				Optional:            true,
				Sensitive:           true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "DVLS API key `$DVLS_API_KEY`. Conflicts with `app_id` / `app_secret`.",
				Optional:            true,
				Sensitive:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("app_id"), path.MatchRoot("app_secret")),
				},
			},
		},
	}
}

func (p *DvlsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data DvlsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	baseUri := data.BaseUri.ValueString()
	appId := stringOrEnv(data.AppId, "DVLS_APP_ID")
	appSecret := stringOrEnv(data.AppSecret, "DVLS_APP_SECRET")
	envApiKey := os.Getenv("DVLS_API_KEY")

	var dvlsClient dvls.Client
	var err error

	switch {
	case !data.ApiKey.IsNull():
		dvlsClient, err = dvls.NewClientWithApiKey(data.ApiKey.ValueString(), baseUri)
	case appId != "" && appSecret != "":
		dvlsClient, err = dvls.NewClient(appId, appSecret, baseUri)
	case envApiKey != "":
		dvlsClient, err = dvls.NewClientWithApiKey(envApiKey, baseUri)
	default:
		resp.Diagnostics.AddError("unable to set up dvls client", "either 'api_key' or both 'app_id' and 'app_secret' must be set")
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("unable to set up dvls client", err.Error())
		return
	}

	resp.DataSourceData = &dvlsClient
	resp.ResourceData = &dvlsClient
	resp.EphemeralResourceData = &dvlsClient
}

func (p *DvlsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewEntryCertificateResource,
		NewEntryCredentialApiKeyResource,
		NewEntryCredentialAzureServicePrincipalResource,
		NewEntryCredentialConnectionStringResource,
		NewEntryCredentialSecretResource,
		NewEntryCredentialSSHKeyResource,
		NewEntryCredentialUsernamePasswordResource,
		NewEntryFolderResource,
		NewVaultResource,
	}
}

func (p *DvlsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewEntryCertificateDataSource,
		NewEntryCredentialApiKeyDataSource,
		NewEntryCredentialAzureServicePrincipalDataSource,
		NewEntryCredentialConnectionStringDataSource,
		NewEntryCredentialSecretDataSource,
		NewEntryCredentialSSHKeyDataSource,
		NewEntryCredentialUsernamePasswordDataSource,
		NewEntryHostDataSource,
		NewEntryWebsiteDataSource,
		NewVaultDataSource,
	}
}

func (p *DvlsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewEntryCertificateEphemeralResource,
		NewEntryCredentialApiKeyEphemeralResource,
		NewEntryCredentialAzureServicePrincipalEphemeralResource,
		NewEntryCredentialConnectionStringEphemeralResource,
		NewEntryCredentialSecretEphemeralResource,
		NewEntryCredentialSSHKeyEphemeralResource,
		NewEntryCredentialUsernamePasswordEphemeralResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &DvlsProvider{
			version: version,
		}
	}
}

func stringOrEnv(value types.String, env string) string {
	if !value.IsNull() {
		return value.ValueString()
	}

	return os.Getenv(env)
}
