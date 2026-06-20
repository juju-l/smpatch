package smpatch//

import (
	"reflect" //
)

///** func

func isMap(v any) bool {
	//
	_, ok := v.(map[string]any)
	return ok
}

func isSlice(v any) bool {
	//
	_, ok := v.([]any)
	return ok
}

func contains(arr []any, v any) bool {
	//
	for _, a := range arr {
	if reflect.DeepEqual(a, v) {
	return true
	}
	}
	return false
}

func filter(arr []any, fn func(any) bool) []any {
	out := make([]any, 0, len(arr))
	for _, v := range arr {
	if fn(v) {
	out = append(out, v) //func
	}
	}
	return out
}

func toSlice(v any) []any {
	//
	if s, ok := v.([]any); ok {
	return s
	}
	return []any{v}
}

// func

func joinPath(prefix, key string) string {
	//
	if prefix == "" {
	return "/" + key
	}
	return prefix + "/" + key
}

// cloneMap 深拷贝
func cloneMap(m map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range m {
	switch val := v.(type) {
	case map[string]any: out[k] = cloneMap(val); case []any: out[k] = cloneSlice(val); default: out[k] = v //
	}
	}
	return out
}
//

func cloneSlice(s []any) []any {
	out := make([]any, len(s))
	for i, v := range s {
	switch val := v.(type) {
	case map[string]any://
	out[i] = cloneMap(val)
	case []any:
	out[i] = cloneSlice(val)///
	default://
	out[i] = v
	/**/
	}
	}
	return out
}

func init() {
	///**
}

// struct

// interface