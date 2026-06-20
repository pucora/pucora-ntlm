package ntlm

import (
	"context"
	"net/http"
	"testing"

	"github.com/pucora/lura/v2/config"
	"github.com/pucora/lura/v2/transport/http/client"
)

func TestConfigGetterRequiresCredentials(t *testing.T) {
	cfg := &config.Backend{
		ExtraConfig: config.ExtraConfig{
			Namespace: map[string]interface{}{
				"user": "DOMAIN\\svc",
			},
		},
	}
	if _, ok := configGetter(cfg.ExtraConfig); ok {
		t.Fatal("expected missing password to disable config")
	}
}

func TestConfigGetterParsesFields(t *testing.T) {
	cfg := &config.Backend{
		ExtraConfig: config.ExtraConfig{
			Namespace: map[string]interface{}{
				"user":     "DOMAIN\\svc",
				"password": "secret",
			},
		},
	}
	got, ok := configGetter(cfg.ExtraConfig)
	if !ok || got.User != `DOMAIN\svc` || got.Password != "secret" {
		t.Fatalf("unexpected config: %+v ok=%v", got, ok)
	}
}

func TestNewHTTPClientWrapsTransport(t *testing.T) {
	cfg := &config.Backend{
		ExtraConfig: config.ExtraConfig{
			Namespace: map[string]interface{}{
				"user":     "alice",
				"password": "secret",
			},
		},
	}
	baseFactory := client.HTTPClientFactory(func(_ context.Context) *http.Client {
		return &http.Client{}
	})
	factory := NewHTTPClient(cfg, baseFactory)
	httpClient := factory(context.Background())
	if httpClient.Transport == nil {
		t.Fatal("expected wrapped transport")
	}
	if _, ok := httpClient.Transport.(*ntlmRoundTripper); !ok {
		t.Fatalf("expected ntlmRoundTripper, got %T", httpClient.Transport)
	}
}

func TestNewHTTPClientPassthroughWithoutConfig(t *testing.T) {
	cfg := &config.Backend{ExtraConfig: config.ExtraConfig{}}
	base := &http.Client{}
	baseFactory := client.HTTPClientFactory(func(_ context.Context) *http.Client { return base })
	factory := NewHTTPClient(cfg, baseFactory)
	if factory(context.Background()) != base {
		t.Fatal("expected same client instance")
	}
}
