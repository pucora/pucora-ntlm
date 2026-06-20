package ntlm

import (
	"github.com/pucora/lura/v2/config"
)

// Namespace is the key to use to store and access the custom config data.
const Namespace = "github.com/pucora/pucora-ntlm"

// Config holds NTLM client authentication settings.
type Config struct {
	User     string
	Password string
}

func configGetter(e config.ExtraConfig) (Config, bool) {
	v, ok := e[Namespace]
	if !ok {
		return Config{}, false
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return Config{}, false
	}
	cfg := Config{}
	if v, ok := tmp["user"]; ok {
		cfg.User, _ = v.(string)
	}
	if v, ok := tmp["password"]; ok {
		cfg.Password, _ = v.(string)
	}
	if cfg.User == "" || cfg.Password == "" {
		return Config{}, false
	}
	return cfg, true
}
