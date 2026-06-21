package smpatch

import (
	"fmt"
	"bytes"
	"gopkg.in/yaml.v3"
)

type Patch struct {
	Ope     string `yaml:"ope"`
	PathKey string `yaml:"pathKey"`

	ByKey   string `yaml:"byKey,omitempty"`
	ItemOps string `yaml:"itemOps,omitempty"`
	MixedAr bool   `yaml:"mixedAr,omitempty"`

	Old   any `yaml:"old,omitempty"`
	Value any `yaml:"value"`
}

func normalizeMap(v any) any {
	switch val := v.(type) {
	case map[any]any:
		m := map[string]any{}
		for k, vv := range val {
			m[fmt.Sprint(k)] = normalizeMap(vv)
		}
		return m

	case map[string]any:
		m := map[string]any{}
		for k, vv := range val {
			m[k] = normalizeMap(vv)
		}
		return m

	case []any:
		for i, vv := range val {
			val[i] = normalizeMap(vv)
		}
		return val

	default:
		return v
	}
}

// cloneViaYAML 使用 YAML Marshal / Unmarshal 实现深拷贝
// 泛型 T 仅用于约束返回类型，内部实现使用 any
func cloneViaYAML[T any](v any) T {
	var buf bytes.Buffer

	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		panic(fmt.Sprintf("yaml encode failed: %v", err))
	}

	var out T
	dec := yaml.NewDecoder(bytes.NewReader(buf.Bytes()))
	dec.KnownFields(true)

	if err := dec.Decode(&out); err != nil {
		panic(fmt.Sprintf("yaml decode failed: %v", err))
	}

	// ✅ 强制修正 map[any]any → map[string]any
	return normalizeMap(out).(T)
}

func copyMap(src, dst map[string]any) {
	for k, v := range src {
		switch v := v.(type) {
		case map[string]any:
			sub := map[string]any{}
			copyMap(v, sub)
			dst[k] = sub
		case []any:
			dst[k] = cloneSlice(v)
		default:
			dst[k] = v
		}
	}
}

func cloneSlice(s []any) []any {
	out := make([]any, len(s))
	for i, v := range s {
		switch v := v.(type) {
		case map[string]any:
			sub := map[string]any{}
			copyMap(v, sub)
			out[i] = sub
		case []any:
			out[i] = cloneSlice(v)
		default:
			out[i] = v
		}
	}
	return out
}

func copyEntireStructure(src, dst map[string]any) {
	for k, v := range src {
		switch v := v.(type) {
		case map[string]any:
			sub := map[string]any{}
			copyEntireStructure(v, sub)
			dst[k] = sub
		case []any:
			dst[k] = cloneSlice(v)
		default:
			dst[k] = v
		}
	}
}

func cloneSlice1(s []any) []any {
	out := make([]any, len(s))
	for i, v := range s {
		switch v := v.(type) {
		case map[string]any:
			sub := map[string]any{}
			copyEntireStructure(v, sub)
			out[i] = sub
		case []any:
			out[i] = cloneSlice(v)
		default:
			out[i] = v
		}
	}
	return out
}