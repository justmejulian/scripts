package registry

import (
	"fmt"
	"scripts/internal/ai/spec"
)

type ProviderFunc func() (spec.Provider, error)

var providers = map[string]ProviderFunc{}

func Register(name string, fn ProviderFunc) {
	if _, exists := providers[name]; exists {
		panic("ai/registry: provider already registered: " + name)
	}
	providers[name] = fn
}

func New(name string) (spec.Provider, error) {
	fn, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("ai/registry: unknown provider %q", name)
	}
	return fn()
}
