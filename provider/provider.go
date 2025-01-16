package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/v2/provider"
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

func (p *CTFdCMProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ctfdcm"
	resp.Version = p.version
}

func (p *CTFdCMProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewChallengeDynamicIaCResource,
	}
}

func (p *CTFdCMProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
