// Contains a wrapper around github.com/ctfer-io/go-ctfd.
//
// It injects spans for all API operations, with improved consistency on typings based
// upon internal assumptions (e.g. IDs are all strings in TF while CTFd has integers,
// so integers -> string is safe as a the first is a subset of the second, but
// also string -> integer as CTFd dictates IDs while the user has no power over it).

package provider

import (
	"context"
	"net/http"

	ctfd "github.com/ctfer-io/go-ctfd/api"
	ctfdcm "github.com/ctfer-io/go-ctfdcm/api"
	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/v2/provider"
	"github.com/ctfer-io/terraform-provider-ctfd/v2/provider/utils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var apiTransport = ctfd.WithTransport(otelhttp.NewTransport(http.DefaultTransport))

func options(ctx context.Context) []ctfd.Option {
	return []ctfd.Option{
		ctfd.WithContext(ctx),
		apiTransport,
	}
}

func GetNonceAndSession(ctx context.Context, url string) (nonce, session string, err error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfd.GetNonceAndSession(url, options(ctx)...)
}

type Client struct {
	sub *ctfd.Client
}

func NewClient(url, nonce, session, apiKey string) *Client {
	return &Client{
		sub: ctfd.NewClient(url, nonce, session, apiKey),
	}
}

func (cli *Client) Login(ctx context.Context, params *ctfd.LoginParams) error {
	ctx, span := StartAPISpan(ctx)
	defer span.End()

	return cli.sub.Login(params, options(ctx)...)
}

// region challenges

func (cli *Client) GetChallenges(ctx context.Context, params *ctfd.GetChallengesParams) ([]*ctfd.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.GetChallenges(params, options(ctx)...)
}

func (cli *Client) GetChallenge(ctx context.Context, id string) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.GetChallenge(cli.sub, id, options(ctx)...)
}

func (cli *Client) PostChallenges(ctx context.Context, params *ctfdcm.PostChallengesParams) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.PostChallenges(cli.sub, params, options(ctx)...)
}

func (cli *Client) PatchChallenges(ctx context.Context, id string, params *ctfdcm.PatchChallengeParams) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.PatchChallenges(cli.sub, id, params, options(ctx)...)
}

func (cli *Client) DeleteChallenge(ctx context.Context, id string) error {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.DeleteChallenge(utils.Atoi(id), options(ctx)...)
}

func (cli *Client) GetChallengeTags(ctx context.Context, id string) ([]*ctfd.Tag, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.GetChallengeTags(utils.Atoi(id), options(ctx)...)
}

func (cli *Client) GetChallengeTopics(ctx context.Context, id string) ([]*ctfd.Topic, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.GetChallengeTopics(utils.Atoi(id), options(ctx)...)
}

func (cli *Client) GetChallengeRequirements(ctx context.Context, id string) (*ctfd.Requirements, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.GetChallengeRequirements(utils.Atoi(id), options(ctx)...)
}

// region instances

func (cli *Client) GetAdminInstance(ctx context.Context, params *ctfdcm.GetAdminInstanceParams) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.GetAdminInstance(cli.sub, params, options(ctx)...)
}

func (cli *Client) PostAdminInstance(ctx context.Context, params *ctfdcm.PostAdminInstanceParams) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.PostAdminInstance(cli.sub, params, options(ctx)...)
}

func (cli *Client) DeleteAdminInstance(ctx context.Context, params *ctfdcm.DeleteAdminInstanceParams) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return ctfdcm.DeleteAdminInstance(cli.sub, params, options(ctx)...)
}

// region tags

func (cli *Client) PostTags(ctx context.Context, params *ctfd.PostTagsParams) (*ctfd.Tag, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.PostTags(params, options(ctx)...)
}

func (cli *Client) DeleteTag(ctx context.Context, id string) error {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.DeleteTag(id, options(ctx)...)
}

// region topics

func (cli *Client) PostTopics(ctx context.Context, params *ctfd.PostTopicsParams) (*ctfd.Topic, error) {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.PostTopics(params, options(ctx)...)
}

func (cli *Client) DeleteTopic(ctx context.Context, params *ctfd.DeleteTopicArgs) error {
	ctx, span := tfctfd.StartAPISpan(ctx)
	defer span.End()

	return cli.sub.DeleteTopic(params, options(ctx)...)
}
