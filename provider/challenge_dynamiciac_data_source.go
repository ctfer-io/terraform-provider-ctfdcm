package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/ctfer-io/go-ctfd/api"
	"github.com/ctfer-io/terraform-provider-ctfd/v2/provider/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = (*challengeDynamicIaCDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*challengeDynamicIaCDataSource)(nil)
)

func NewChallengeDynamicIaCDataSource() datasource.DataSource {
	return &challengeDynamicIaCDataSource{}
}

type challengeDynamicIaCDataSource struct {
	client *api.Client
}

type challengesDynamicDataSourceModel struct {
	ID         types.String                       `tfsdk:"id"`
	Challenges []ChallengeDynamicIaCResourceModel `tfsdk:"challenges"`
}

func (ch *challengeDynamicIaCDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_challenges_dynamiciac"
}

func (ch *challengeDynamicIaCDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"challenges": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Identifier of the challenge.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the challenge, displayed as it.",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "Category of the challenge that CTFd groups by on the web UI.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the challenge, consider using multiline descriptions for better style.",
							Computed:            true,
						},
						"attribution": schema.StringAttribute{
							MarkdownDescription: "Attribution to the creator(s) of the challenge.",
							Computed:            true,
						},
						"connection_info": schema.StringAttribute{
							MarkdownDescription: "Connection Information to connect to the challenge instance, useful for pwn or web pentest.",
							Computed:            true,
						},
						"max_attempts": schema.Int64Attribute{
							MarkdownDescription: "Maximum amount of attempts before being unable to flag the challenge.",
							Computed:            true,
						},
						"value": schema.Int64Attribute{
							MarkdownDescription: "The value (points) of the challenge once solved. It is mapped to `initial` under the hood, but displayed as `value` for consistency with the standard challenge.",
							Computed:            true,
						},
						"decay": schema.Int64Attribute{
							MarkdownDescription: "The decay defines from each number of solves does the decay function triggers until reaching minimum. This function is defined by CTFd and could be configured through `.function`.",
							Computed:            true,
						},
						"minimum": schema.Int64Attribute{
							MarkdownDescription: "The minimum points for a dynamic-score challenge to reach with the decay function. Once there, no solve could have more value.",
							Computed:            true,
						},
						"function": schema.StringAttribute{
							MarkdownDescription: "Decay function to define how the challenge value evolve through solves, either linear or logarithmic.",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "State of the challenge, either hidden or visible.",
							Computed:            true,
						},
						"next": schema.Int64Attribute{
							MarkdownDescription: "Suggestion for the end-user as next challenge to work on.",
							Computed:            true,
						},
						"requirements": schema.SingleNestedAttribute{
							MarkdownDescription: "List of required challenges that needs to get flagged before this one being accessible. Useful for skill-trees-like strategy CTF.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"behavior": schema.StringAttribute{
									MarkdownDescription: "Behavior if not unlocked, either hidden or anonymized.",
									Computed:            true,
								},
								"prerequisites": schema.ListAttribute{
									MarkdownDescription: "List of the challenges ID.",
									Computed:            true,
									ElementType:         types.StringType,
								},
							},
						},
						"tags": schema.ListAttribute{
							MarkdownDescription: "List of challenge tags that will be displayed to the end-user. You could use them to give some quick insights of what a challenge involves.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"topics": schema.ListAttribute{
							MarkdownDescription: "List of challenge topics that are displayed to the administrators for maintenance and planification.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"shared": schema.BoolAttribute{
							MarkdownDescription: "Whether the instance will be shared between all players.",
							Computed:            true,
						},
						"destroy_on_flag": schema.BoolAttribute{
							MarkdownDescription: "Whether to destroy the instance once flagged.",
							Computed:            true,
						},
						"mana_cost": schema.Int64Attribute{
							MarkdownDescription: "The cost (in mana) of the challenge once an instance is deployed.",
							Computed:            true,
						},
						"scenario": schema.StringAttribute{
							MarkdownDescription: "The OCI reference to the scenario.",
							Computed:            true,
						},
						"timeout": schema.Int64Attribute{
							MarkdownDescription: "The timeout (in seconds) after which the instance will be janitored.",
							Computed:            true,
						},
						"until": schema.StringAttribute{
							MarkdownDescription: "The date until the instance could run before being janitored.",
							Computed:            true,
						},
						"additional": schema.MapAttribute{
							MarkdownDescription: "An optional key=value map (both strings) to pass to the scenario.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"min": schema.Int64Attribute{
							MarkdownDescription: "The minimum number of instances to set in the pool.",
							Computed:            true,
						},
						"max": schema.Int64Attribute{
							MarkdownDescription: "The number of instances after which not to pool anymore.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (ch *challengeDynamicIaCDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *github.com/ctfer-io/go-ctfd/api.Client, got: %T. Please open an issue at https://github.com/ctfer-io/terraform-provider-ctfdcm", req.ProviderData),
		)
		return
	}

	ch.client = client
}

func (ch *challengeDynamicIaCDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state challengesDynamicDataSourceModel

	challs, err := ch.client.GetChallenges(&api.GetChallengesParams{
		Type: utils.Ptr("dynamic_iac"),
	}, api.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read CTFd Challenges",
			err.Error(),
		)
		return
	}

	state.Challenges = make([]ChallengeDynamicIaCResourceModel, 0, len(challs))
	for _, c := range challs {
		chall := ChallengeDynamicIaCResourceModel{}
		chall.ID = types.StringValue(strconv.Itoa(c.ID))
		chall.Read(ctx, ch.client, resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		state.Challenges = append(state.Challenges, chall)
	}

	state.ID = types.StringValue("placeholder")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
