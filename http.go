package ntlm

import (
	"context"
	"net/http"

	"github.com/Azure/go-ntlmssp"
	"github.com/pucora/lura/v2/config"
	"github.com/pucora/lura/v2/transport/http/client"
)

// NewHTTPClient wraps an HTTP client factory with NTLMv2 transport when configured.
func NewHTTPClient(cfg *config.Backend, next client.HTTPClientFactory) client.HTTPClientFactory {
	ntlmCfg, ok := configGetter(cfg.ExtraConfig)
	if !ok {
		return next
	}
	user := ntlmCfg.User
	password := ntlmCfg.Password
	return func(ctx context.Context) *http.Client {
		base := next(ctx)
		transport := base.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		base.Transport = &ntlmRoundTripper{
			user:     user,
			password: password,
			next: ntlmssp.Negotiator{
				RoundTripper: transport,
			},
		}
		return base
	}
}

type ntlmRoundTripper struct {
	user     string
	password string
	next     http.RoundTripper
}

func (rt *ntlmRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	cloned.SetBasicAuth(rt.user, rt.password)
	return rt.next.RoundTrip(cloned)
}
