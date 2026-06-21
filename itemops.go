package smpatch

import (
	"strings"
	"slices"
)

func itemOps(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			cur[parts[i]] = map[string]any{}
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	arr := DeepCopy(cur[key]).([]any)
	vals := p.Value.([]any)

	switch p.ItemOps {
	case "add":
		for _, v := range vals {
			if !slices.Contains(arr, v) {
				arr = append(arr, v)
			}
		}
	case "remove":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(vals, v)
		})
	case "replace":
		for i, v := range arr {
			if v == p.Old {
				arr[i] = vals[0]
				break
			}
		}
	case "keep":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return !slices.Contains(vals, v)
		})
	case "disable":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(vals, v)
		})
	}

	cur[key] = arr
	return nil
}