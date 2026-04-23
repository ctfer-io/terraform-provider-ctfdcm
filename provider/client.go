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

func apiOptions(ctx context.Context) []ctfd.Option {
	return []ctfd.Option{
		ctfd.WithContext(ctx),
		apiTransport,
	}
}

func GetNonceAndSession(ctx context.Context, url string, opts ...Option) (nonce, session string, err error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfd.GetNonceAndSession(url, apiOptions(ctx)...)
}

type Client struct {
	sub *ctfd.Client
}

func NewClient(url, nonce, session, apiKey string) *Client {
	return &Client{
		sub: ctfd.NewClient(url, nonce, session, apiKey),
	}
}

func (cli *Client) Login(ctx context.Context, params *ctfd.LoginParams, opts ...Option) error {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.Login(params, apiOptions(ctx)...)
}

// region challenges

func (cli *Client) GetChallenges(ctx context.Context, params *ctfd.GetChallengesParams, opts ...Option) ([]*ctfd.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.GetChallenges(params, apiOptions(ctx)...)
}

func (cli *Client) GetChallenge(ctx context.Context, id string, opts ...Option) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.GetChallenge(cli.sub, id, apiOptions(ctx)...)
}

func (cli *Client) PostChallenges(ctx context.Context, params *ctfdcm.PostChallengesParams, opts ...Option) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.PostChallenges(cli.sub, params, apiOptions(ctx)...)
}

func (cli *Client) PatchChallenges(ctx context.Context, id string, params *ctfdcm.PatchChallengeParams, opts ...Option) (*ctfdcm.Challenge, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.PatchChallenges(cli.sub, id, params, apiOptions(ctx)...)
}

func (cli *Client) DeleteChallenge(ctx context.Context, id string, opts ...Option) error {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.DeleteChallenge(utils.Atoi(id), apiOptions(ctx)...)
}

func (cli *Client) GetChallengeTags(ctx context.Context, id string, opts ...Option) ([]*ctfd.Tag, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.GetChallengeTags(utils.Atoi(id), apiOptions(ctx)...)
}

func (cli *Client) GetChallengeTopics(ctx context.Context, id string, opts ...Option) ([]*ctfd.Topic, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.GetChallengeTopics(utils.Atoi(id), apiOptions(ctx)...)
}

func (cli *Client) GetChallengeRequirements(ctx context.Context, id string, opts ...Option) (*ctfd.Requirements, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.GetChallengeRequirements(utils.Atoi(id), apiOptions(ctx)...)
}

// region instances

func (cli *Client) GetAdminInstance(ctx context.Context, params *ctfdcm.GetAdminInstanceParams, opts ...Option) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.GetAdminInstance(cli.sub, params, apiOptions(ctx)...)
}

func (cli *Client) PostAdminInstance(ctx context.Context, params *ctfdcm.PostAdminInstanceParams, opts ...Option) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.PostAdminInstance(cli.sub, params, apiOptions(ctx)...)
}

func (cli *Client) DeleteAdminInstance(ctx context.Context, params *ctfdcm.DeleteAdminInstanceParams, opts ...Option) (*ctfdcm.Instance, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return ctfdcm.DeleteAdminInstance(cli.sub, params, apiOptions(ctx)...)
}

// region tags

func (cli *Client) PostTags(ctx context.Context, params *ctfd.PostTagsParams, opts ...Option) (*ctfd.Tag, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.PostTags(params, apiOptions(ctx)...)
}

func (cli *Client) DeleteTag(ctx context.Context, id string, opts ...Option) error {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.DeleteTag(id, apiOptions(ctx)...)
}

// region topics

func (cli *Client) PostTopics(ctx context.Context, params *ctfd.PostTopicsParams, opts ...Option) (*ctfd.Topic, error) {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.PostTopics(params, apiOptions(ctx)...)
}

func (cli *Client) DeleteTopic(ctx context.Context, params *ctfd.DeleteTopicArgs, opts ...Option) error {
	ctx, span := tfctfd.StartAPISpan(ctx, getTracer(opts...))
	defer span.End()

	return cli.sub.DeleteTopic(params, apiOptions(ctx)...)
}
