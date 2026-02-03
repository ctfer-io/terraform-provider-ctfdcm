package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	ctfdcm "github.com/ctfer-io/go-ctfdcm/api"
)

var (
	_ resource.Resource              = (*instanceResource)(nil)
	_ resource.ResourceWithConfigure = (*instanceResource)(nil)
)

func NewInstanceResource() resource.Resource {
	return &instanceResource{}
}

type instanceResource struct {
	client *Client
}

type InstanceResourceModel struct {
	ChallengeID types.String `tfsdk:"challenge_id"`
	SourceID    types.String `tfsdk:"source_id"`
}

func (r *instanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_instance"
}

func (r *instanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "CTFd is built around the Challenge resource, which contains all the attributes to define a part of the Capture The Flag event.\n\nThis implementation has support of On Demand infrastructures through [Chall-Manager](https://github.com/ctfer-io/chall-manager).",
		Attributes: map[string]schema.Attribute{
			"challenge_id": schema.StringAttribute{
				MarkdownDescription: "The challenge to provision an instance of.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "The source of whom to provision an instance for.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *instanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected %T, got: %T. Please open an issue at https://github.com/ctfer-io/terraform-provider-ctfdcm", (*Client)(nil), req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *instanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, span := StartTFSpan(ctx, r)
	defer span.End()

	var data InstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.PostAdminInstance(ctx, &ctfdcm.PostAdminInstanceParams{
		ChallengeID: data.ChallengeID.ValueString(),
		SourceID:    data.SourceID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create instance, got error: %s", err),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *instanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	ctx, span := StartTFSpan(ctx, r)
	defer span.End()

	var data InstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.GetAdminInstance(ctx, &ctfdcm.GetAdminInstanceParams{
		ChallengeID: data.ChallengeID.ValueString(),
		SourceID:    data.SourceID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read instance, got error: %s", err),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *instanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// It should not happen

	ctx, span := StartTFSpan(ctx, r)
	defer span.End()

	var data InstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Provider Error", "CTFd-Chall-Manager does not permit update of instance-related information thus this provider cannot do so. This operation should not have been possible.")
}

func (r *instanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	ctx, span := StartTFSpan(ctx, r)
	defer span.End()

	var data InstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.DeleteAdminInstance(ctx, &ctfdcm.DeleteAdminInstanceParams{
		ChallengeID: data.ChallengeID.ValueString(),
		SourceID:    data.SourceID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete instance, got error: %s", err),
		)
		return
	}
}
