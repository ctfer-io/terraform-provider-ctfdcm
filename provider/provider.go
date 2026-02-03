package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/ctfer-io/go-ctfd/api"
	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/v2/provider"
	"github.com/ctfer-io/terraform-provider-ctfd/v2/provider/utils"
)

const (
	providerTypeName = "ctfdcm"
)

var _ provider.Provider = (*CTFdCMProvider)(nil)

type CTFdCMProvider struct {
	version string
	*tfctfd.CTFdProvider
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CTFdCMProvider{
			version: version,
		}
	}
}

type CTFdCMProviderModel struct {
	URL      types.String `tfsdk:"url"`
	APIKey   types.String `tfsdk:"api_key"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *CTFdCMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = providerTypeName
	resp.Version = p.version
}

func (p *CTFdCMProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CTFdCMProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check configuration values are known
	if config.URL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("url"),
			"Unknown CTFD url.",
			"The provider cannot guess where to reach the CTFd instance.",
		)
	}
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown CTFd API key.",
			"The provider cannot create the CTFd API client as there is an unknown API key value.",
		)
	}
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown CTFd admin or service account username.",
			"The provider cannot create the CTFd API client as there is an unknown username.",
		)
	}
	if config.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown CTFd admin or service account password.",
			"The provider cannot create the CTFd API client as there is an unknown password.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract environment variables values
	url := os.Getenv("CTFD_URL")
	apiKey := os.Getenv("CTFD_API_KEY")
	username := os.Getenv("CTFD_ADMIN_USERNAME")
	password := os.Getenv("CTFD_ADMIN_PASSWORD")

	if !config.URL.IsNull() {
		url = config.URL.ValueString()
	}
	if !config.APIKey.IsNull() {
		apiKey = config.APIKey.ValueString()
	}
	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}
	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// Check there is enough content
	ak := apiKey != ""
	up := username != "" && password != ""
	if !ak && !up {
		resp.Diagnostics.AddError(
			"CTFd provider configuration error",
			"The provider cannot create the CTFd API client as there is an invalid configuration. Expected either an API key, a nonce and session, or a username and password.",
		)
		return
	}

	// Instantiate CTFd API client
	ctx = tflog.SetField(ctx, "ctfd_url", url)
	ctx = utils.AddSensitive(ctx, "ctfd_api_key", apiKey)
	ctx = utils.AddSensitive(ctx, "ctfd_username", username)
	ctx = utils.AddSensitive(ctx, "ctfd_password", password)
	tflog.Debug(ctx, "Creating CTFd API client")

	nonce, session, err := GetNonceAndSession(ctx, url)
	if err != nil {
		resp.Diagnostics.AddError(
			"CTFd error",
			fmt.Sprintf("Failed to fetch nonce and session: %s", err),
		)
		return
	}

	client := NewClient(url, nonce, session, apiKey)
	if up {
		// XXX due to the CTFd ratelimiter on rare endpoint
		if _, ok := os.LookupEnv("TF_ACC"); ok {
			time.Sleep(5 * time.Second)
		}

		if err := client.Login(ctx, &api.LoginParams{
			Name:     username,
			Password: password,
		}); err != nil {
			resp.Diagnostics.AddError(
				"CTFd error",
				fmt.Sprintf("Failed to login: %s", err),
			)
			return
		}
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configure CTFd API client", map[string]any{
		"success": true,
		"login":   up,
	})
}

func (p *CTFdCMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewChallengeDynamicIaCResource,
		NewInstanceResource,
	}
}

func (p *CTFdCMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewChallengeDynamicIaCDataSource,
	}
}
