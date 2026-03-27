package requestconfig

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseJSON(provider, raw string) (map[string]any, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	var value any
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		return nil, fmt.Errorf("%s: invalid request config: %w", provider, err)
	}

	values, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s: invalid request config: must be a JSON object", provider)
	}

	return values, nil
}

func EnsureNoOverrides(provider string, values map[string]any, forbidden ...string) error {
	blocked := make(map[string]struct{}, len(forbidden))
	for _, key := range forbidden {
		blocked[key] = struct{}{}
	}

	for key := range values {
		if _, found := blocked[key]; found {
			return fmt.Errorf("%s: request config may not override %q", provider, key)
		}
	}

	return nil
}

func Merge(base, extra map[string]any) map[string]any {
	merged := make(map[string]any, len(base)+len(extra))

	for key, value := range base {
		merged[key] = value
	}

	for key, value := range extra {
		merged[key] = value
	}

	return merged
}

func Apply(provider string, base map[string]any, raw string, forbidden ...string) (map[string]any, error) {
	extra, err := ParseJSON(provider, raw)
	if err != nil {
		return nil, err
	}

	if err := EnsureNoOverrides(provider, extra, forbidden...); err != nil {
		return nil, err
	}

	return Merge(base, extra), nil
}
