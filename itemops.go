package smpatch

import (
	"fmt"
	"slices"
	"strings"
)

func itemOps(
	p *Patch,
	src map[string]any,
	tgt map[string]any,
) error {

	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := src
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			cur[parts[i]] = map[string]any{}
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	arr := cloneViaYAML[[]any](cur[key])

	// ✅ 统一转成 []any
	vals, ok := p.Value.([]any)
	if !ok {
		return fmt.Errorf("itemOps %s requires value to be array", p.ItemOps)
	}

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
		if p.Old == nil {
			return fmt.Errorf("itemOps replace requires old")
		}
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

	tgt[key] = arr
	return nil
}