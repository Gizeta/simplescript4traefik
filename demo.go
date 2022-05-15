// Package simplescript4traefik a demo plugin.
package simplescript4traefik

import (
	"context"
	"fmt"
	"net/http"
)

// Config the plugin configuration.
type Config struct {
	Code 	string `json:"code"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Code: "",
	}
}

// Demo a ScriptPlugin plugin.
type ScriptPlugin struct {
	name     string
	next     http.Handler
	code     string
}

// New created a new ScriptPlugin plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ScriptPlugin{
		name: name,
		next: next,
		code: config.Code,
	}, nil
}

func (a *ScriptPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if a.code != "" {
		env := CreateEnv()
		RegisterTraefikBuiltin(&env, &rw, req, &(a.next))
		RunScript(a.code, &env)
		return
	}

	a.next.ServeHTTP(rw, req)
}
