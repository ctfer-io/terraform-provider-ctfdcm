package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	ctfd "github.com/ctfer-io/go-ctfd/api"
	ctfdcm "github.com/ctfer-io/go-ctfdcm/api"
	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/v2/provider"
	"github.com/ctfer-io/terraform-provider-ctfd/v2/provider/utils"
)

var (
	_ resource.Resource = (*challengeDynamicIaCResource)(nil)
	_ resource.Resource = (*challengeDynamicIaCResource)(nil)
	_ resource.Resource = (*challengeDynamicIaCResource)(nil)
)

func NewChallengeDynamicIaCResource() resource.Resource {
	return &challengeDynamicIaCResource{}
}

type challengeDynamicIaCResource struct {
	client *ctfd.Client
}

type ChallengeDynamicIaCResourceModel struct {
	tfctfd.ChallengeDynamicResourceModel

	Shared        types.Bool   `tfsdk:"shared"`
	DestroyOnFlag types.Bool   `tfsdk:"destroy_on_flag"`
	ManaCost      types.Int64  `tfsdk:"mana_cost"`
	ScenarioID    types.String `tfsdk:"scenario_id"`
	Timeout       types.Int64  `tfsdk:"timeout"`
	Until         types.String `tfsdk:"until"`
	Additional    types.Map    `tfsdk:"additional"`
	Min           types.Int64  `tfsdk:"min"`
	Max           types.Int64  `tfsdk:"max"`
}

func (r *challengeDynamicIaCResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_challenge_dynamiciac"
}

func (r *challengeDynamicIaCResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "CTFd is built around the Challenge resource, which contains all the attributes to define a part of the Capture The Flag event.\n\nThis implementation has support of On Demand infrastructures through [Chall-Manager](https://github.com/ctfer-io/chall-manager).",
		Attributes:          ChallengeDynamicIaCResourceAttributes,
	}
}

func (r *challengeDynamicIaCResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ctfd.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *github.com/ctfer-io/go-ctfd/api.Client, got: %T. Please open an issue at https://github.com/ctfer-io/terraform-provider-ctfdcm", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *challengeDynamicIaCResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ChallengeDynamicIaCResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create Challenge
	reqs := (*ctfd.Requirements)(nil)
	if data.Requirements != nil {
		preqs := make([]int, 0, len(data.Requirements.Prerequisites))
		for _, preq := range data.Requirements.Prerequisites {
			id, _ := strconv.Atoi(preq.ValueString())
			preqs = append(preqs, id)
		}
		reqs = &ctfd.Requirements{
			Anonymize:     tfctfd.GetAnon(data.Requirements.Behavior),
			Prerequisites: preqs,
		}
	}
	add := map[string]string{}
	for k, tv := range data.Additional.Elements() {
		add[k] = tv.(types.String).ValueString()
	}
	res, err := ctfdcm.PostChallenges(r.client, &ctfdcm.PostChallengesParams{
		// CTFd
		Name:           data.Name.ValueString(),
		Category:       data.Category.ValueString(),
		Description:    data.Description.ValueString(),
		Attribution:    data.Attribution.ValueStringPointer(),
		ConnectionInfo: data.ConnectionInfo.ValueStringPointer(),
		MaxAttempts:    utils.ToInt(data.MaxAttempts),
		Function:       data.Function.ValueStringPointer(),
		Initial:        utils.ToInt(data.Value),
		Decay:          utils.ToInt(data.Decay),
		Minimum:        utils.ToInt(data.Minimum),
		State:          data.State.ValueString(),
		Type:           "dynamic_iac",
		NextID:         utils.ToInt(data.Next),
		Requirements:   reqs,
		// CTFd-Chall-Manager plugin
		DestroyOnFlag: data.DestroyOnFlag.ValueBool(),
		Shared:        data.Shared.ValueBool(),
		ManaCost:      int(data.ManaCost.ValueInt64()),
		ScenarioID:    utils.Atoi(data.ScenarioID.ValueString()),
		Timeout:       utils.ToInt(data.Timeout),
		Until:         data.Until.ValueStringPointer(),
		Additional:    add,
		Min:           int(data.Min.ValueInt64()),
		Max:           int(data.Max.ValueInt64()),
	}, ctfd.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create challenge, got error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "created a challenge")

	// Save computed attributes in state
	data.ID = types.StringValue(strconv.Itoa(res.ID))
	data.Timeout = utils.ToTFInt64(res.Timeout)
	data.Until = types.StringPointerValue(res.Until)

	// Create tags
	challTags := make([]types.String, 0, len(data.Tags))
	for _, tag := range data.Tags {
		_, err := r.client.PostTags(&ctfd.PostTagsParams{
			Challenge: utils.Atoi(data.ID.ValueString()),
			Value:     tag.ValueString(),
		}, ctfd.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to create tags, got error: %s", err),
			)
			return
		}
		challTags = append(challTags, tag)
	}
	if data.Tags != nil {
		data.Tags = challTags
	}

	// Create topics
	challTopics := make([]types.String, 0, len(data.Topics))
	for _, topic := range data.Topics {
		_, err := r.client.PostTopics(&ctfd.PostTopicsParams{
			Challenge: utils.Atoi(data.ID.ValueString()),
			Type:      "challenge",
			Value:     topic.ValueString(),
		}, ctfd.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to create topic, got error: %s", err),
			)
			return
		}
		challTopics = append(challTopics, topic)
	}
	if data.Topics != nil {
		data.Topics = challTopics
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *challengeDynamicIaCResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ChallengeDynamicIaCResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Read(ctx, r.client, resp.Diagnostics)

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *challengeDynamicIaCResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ChallengeDynamicIaCResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var dataState ChallengeDynamicIaCResourceModel
	req.State.Get(ctx, &dataState)

	// Patch direct attributes
	reqs := (*ctfd.Requirements)(nil)
	if data.Requirements != nil {
		preqs := make([]int, 0, len(data.Requirements.Prerequisites))
		for _, preq := range data.Requirements.Prerequisites {
			id, _ := strconv.Atoi(preq.ValueString())
			preqs = append(preqs, id)
		}
		reqs = &ctfd.Requirements{
			Anonymize:     tfctfd.GetAnon(data.Requirements.Behavior),
			Prerequisites: preqs,
		}
	}
	add := map[string]string{}
	for k, tv := range data.Additional.Elements() {
		add[k] = tv.(types.String).ValueString()
	}
	res, err := ctfdcm.PatchChallenges(r.client, data.ID.ValueString(), &ctfdcm.PatchChallengeParams{
		// CTFd
		Name:           data.Name.ValueString(),
		Category:       data.Category.ValueString(),
		Description:    data.Description.ValueString(),
		Attribution:    data.Attribution.ValueStringPointer(),
		ConnectionInfo: data.ConnectionInfo.ValueStringPointer(),
		MaxAttempts:    utils.ToInt(data.MaxAttempts),
		Function:       data.Function.ValueStringPointer(),
		Initial:        utils.ToInt(data.Value),
		Decay:          utils.ToInt(data.Decay),
		Minimum:        utils.ToInt(data.Minimum),
		State:          data.State.ValueString(),
		NextID:         utils.ToInt(data.Next),
		Requirements:   reqs,
		// CTFd-Chall-Manager plugin
		DestroyOnFlag: data.DestroyOnFlag.ValueBool(),
		Shared:        data.Shared.ValueBool(),
		ManaCost:      int(data.ManaCost.ValueInt64()),
		ScenarioID:    utils.Atoi(data.ScenarioID.ValueString()),
		Timeout:       utils.ToInt(data.Timeout),
		Until:         data.Until.ValueStringPointer(),
		Additional:    add,
		Min:           int(data.Min.ValueInt64()),
		Max:           int(data.Max.ValueInt64()),
	}, ctfd.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to update challenge, got error: %s", err),
		)
		return
	}
	data.Timeout = utils.ToTFInt64(res.Timeout)
	data.Until = types.StringPointerValue(res.Until)

	// Update its tags (drop them all, create new ones)
	challTags, err := r.client.GetChallengeTags(utils.Atoi(data.ID.ValueString()), ctfd.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to get all tags of challenge %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}
	for _, tag := range challTags {
		if err := r.client.DeleteTag(strconv.Itoa(tag.ID), ctfd.WithContext(ctx)); err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to delete tag %d of challenge %s, got error: %s", tag.ID, data.ID.ValueString(), err),
			)
			return
		}
	}
	tags := make([]types.String, 0, len(data.Tags))
	for _, tag := range data.Tags {
		_, err := r.client.PostTags(&ctfd.PostTagsParams{
			Challenge: utils.Atoi(data.ID.ValueString()),
			Value:     tag.ValueString(),
		}, ctfd.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to create tag of challenge %s, got error: %s", data.ID.ValueString(), err),
			)
			return
		}
		tags = append(tags, tag)
	}
	if data.Tags != nil {
		data.Tags = tags
	}

	// Update its topics (drop them all, create new ones)
	challTopics, err := r.client.GetChallengeTopics(utils.Atoi(data.ID.ValueString()), ctfd.WithContext(ctx))
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to get all topics of challenge %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}
	for _, topic := range challTopics {
		if err := r.client.DeleteTopic(&ctfd.DeleteTopicArgs{
			ID:   strconv.Itoa(topic.ID),
			Type: "challenge",
		}, ctfd.WithContext(ctx)); err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to delete topic %d of challenge %s, got error: %s", topic.ID, data.ID.ValueString(), err),
			)
			return
		}
	}
	topics := make([]types.String, 0, len(data.Topics))
	for _, topic := range data.Topics {
		_, err := r.client.PostTopics(&ctfd.PostTopicsParams{
			Challenge: utils.Atoi(data.ID.ValueString()),
			Type:      "challenge",
			Value:     topic.ValueString(),
		}, ctfd.WithContext(ctx))
		if err != nil {
			resp.Diagnostics.AddError(
				"Client Error",
				fmt.Sprintf("Unable to create topic of challenge %s, got error: %s", data.ID.ValueString(), err),
			)
			return
		}
		topics = append(topics, topic)
	}
	if data.Topics != nil {
		data.Topics = topics
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *challengeDynamicIaCResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ChallengeDynamicIaCResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteChallenge(utils.Atoi(data.ID.ValueString()), ctfd.WithContext(ctx)); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete challenge, got error: %s", err))
		return
	}

	// ... don't need to delete nested objects, this is handled by CTFd
}

func (r *challengeDynamicIaCResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// Automatically call r.Read
}

func (chall *ChallengeDynamicIaCResourceModel) Read(ctx context.Context, client *ctfd.Client, diags diag.Diagnostics) {
	res, err := ctfdcm.GetChallenge(client, chall.ID.ValueString(), ctfd.WithContext(ctx))
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Unable to read challenge %s, got error: %s", chall.ID.ValueString(), err))
		return
	}
	// CTFd
	chall.Name = types.StringValue(res.Name)
	chall.Category = types.StringValue(res.Category)
	chall.Description = types.StringValue(res.Description)
	chall.Attribution = types.StringPointerValue(res.Attribution)
	chall.ConnectionInfo = utils.ToTFString(res.ConnectionInfo)
	chall.MaxAttempts = utils.ToTFInt64(res.MaxAttempts)
	chall.Function = utils.ToTFString(res.Function)
	chall.Value = utils.ToTFInt64(res.Initial)
	chall.Decay = utils.ToTFInt64(res.Decay)
	chall.Minimum = utils.ToTFInt64(res.Minimum)
	chall.State = types.StringValue(res.State)
	chall.Next = utils.ToTFInt64(res.NextID)
	// CTFer.io Chall-Manager plugin
	chall.DestroyOnFlag = types.BoolValue(res.DestroyOnFlag)
	chall.Shared = types.BoolValue(res.Shared)
	chall.ManaCost = types.Int64Value(int64(res.ManaCost))
	chall.ScenarioID = types.StringValue(strconv.Itoa(res.ScenarioID))
	chall.Timeout = utils.ToTFInt64(res.Timeout)
	chall.Until = types.StringPointerValue(res.Until)
	addMp := map[string]attr.Value{}
	for k, v := range res.Additional {
		addMp[k] = types.StringValue(v)
	}
	add, d := types.MapValue(types.StringType, addMp)
	diags.Append(d...)
	chall.Additional = add
	chall.Min = types.Int64Value(int64(res.Min))
	chall.Max = types.Int64Value(int64(res.Max))

	id := utils.Atoi(chall.ID.ValueString())

	// Get subresources
	// => Requirements
	resReqs, err := client.GetChallengeRequirements(id, ctfd.WithContext(ctx))
	if err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read challenge %d requirements, got error: %s", id, err),
		)
		return
	}
	reqs := (*tfctfd.RequirementsSubresourceModel)(nil)
	if resReqs != nil {
		challPreqs := make([]types.String, 0, len(resReqs.Prerequisites))
		for _, req := range resReqs.Prerequisites {
			challPreqs = append(challPreqs, types.StringValue(strconv.Itoa(req)))
		}
		reqs = &tfctfd.RequirementsSubresourceModel{
			Behavior:      tfctfd.FromAnon(resReqs.Anonymize),
			Prerequisites: challPreqs,
		}
	}
	chall.Requirements = reqs

	// => Tags
	resTags, err := client.GetChallengeTags(id, ctfd.WithContext(ctx))
	if err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read challenge %d tags, got error: %s", id, err),
		)
		return
	}
	chall.Tags = make([]basetypes.StringValue, 0, len(resTags))
	for _, tag := range resTags {
		chall.Tags = append(chall.Tags, types.StringValue(tag.Value))
	}

	// => Topics
	resTopics, err := client.GetChallengeTopics(id, ctfd.WithContext(ctx))
	if err != nil {
		diags.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read challenge %d topics, got error: %s", id, err),
		)
		return
	}
	chall.Topics = make([]basetypes.StringValue, 0, len(resTopics))
	for _, topic := range resTopics {
		chall.Topics = append(chall.Topics, types.StringValue(topic.Value))
	}
}

var (
	ChallengeDynamicIaCResourceAttributes = utils.BlindMerge(tfctfd.ChallengeDynamicResourceAttributes, map[string]schema.Attribute{
		"shared": schema.BoolAttribute{
			MarkdownDescription: "Whether the instance will be shared between all players.",
			Optional:            true,
			Computed:            true,
			Default:             defaults.Bool(booldefault.StaticBool(false)),
		},
		"destroy_on_flag": schema.BoolAttribute{
			MarkdownDescription: "Whether to destroy the instance once flagged.",
			Optional:            true,
			Computed:            true,
			Default:             defaults.Bool(booldefault.StaticBool(false)),
		},
		"mana_cost": schema.Int64Attribute{
			MarkdownDescription: "The cost (in mana) of the challenge once an instance is deployed.",
			Optional:            true,
			Computed:            true,
			Default:             defaults.Int64(int64default.StaticInt64(0)),
		},
		"scenario_id": schema.StringAttribute{
			MarkdownDescription: "The file's ID of the scenario.",
			Required:            true,
		},
		"timeout": schema.Int64Attribute{
			MarkdownDescription: "The timeout (in seconds) after which the instance will be janitored.",
			Optional:            true,
			Computed:            true,
			Default:             nil,
		},
		"until": schema.StringAttribute{
			MarkdownDescription: "The date until the instance could run before being janitored.",
			Optional:            true,
			Computed:            true,
			Default:             nil,
		},
		"additional": schema.MapAttribute{
			MarkdownDescription: "An optional key=value map (both strings) to pass to the scenario.",
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			Default:             defaults.Map(mapdefault.StaticValue(basetypes.NewMapValueMust(types.StringType, map[string]attr.Value{}))),
		},
		"min": schema.Int64Attribute{
			MarkdownDescription: "The minimum number of instances to set in the pool.",
			Optional:            true,
			Computed:            true,
			Default:             defaults.Int64(int64default.StaticInt64(0)),
		},
		"max": schema.Int64Attribute{
			MarkdownDescription: "The number of instances after which not to pool anymore.",
			Optional:            true,
			Computed:            true,
			Default:             defaults.Int64(int64default.StaticInt64(0)),
		},
	})
)
