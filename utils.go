package yamlpatch

import (
	"reflect"
	"strings"
)

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

func isSlice(v interface{}) bool {
	_, ok := v.([]interface{})
	return ok
}

func contains(arr []interface{}, v interface{}) bool {
	for _, a := range arr {
		if reflect.DeepEqual(a, v) {
			return true
		}
	}
	return false
}

func filter(arr []interface{}, fn func(interface{}) bool) []interface{} {
	out := make([]interface{}, 0, len(arr))
	for _, v := range arr {
		if fn(v) {
			out = append(out, v)
		}
	}
	return out
}

func toSlice(v interface{}) []interface{} {
	if s, ok := v.([]interface{}); ok {
		return s
	}
	return []interface{}{v}
}

func joinPath(prefix, key string) string {
	if prefix == "" {
		return "/" + key
	}
	return prefix + "/" + key
}

// cloneMap 深拷贝
func cloneMap(m map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	for k, v := range m {
		switch val := v.(type) {
		case map[string]interface{}:
			out[k] = cloneMap(val)
		case []interface{}:
			out[k] = cloneSlice(val)
		default:
			out[k] = v
		}
	}
	return out
}

func cloneSlice(s []interface{}) []interface{} {
	out := make([]interface{}, len(s))
	for i, v := range s {
		switch val := v.(type) {
		case map[string]interface{}:
			out[i] = cloneMap(val)
		case []interface{}:
			out[i] = cloneSlice(val)
		default:
			out[i] = v
		}
	}
	return out
}