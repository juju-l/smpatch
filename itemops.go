package smpatch

import (
	"strings"
	"slices"
)

func itemOps(
	p *Patch,
	src map[string]any,
	tgt map[string]any,
) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")
	cur := src
	for i := 0; i < len(parts)-1; i++ {
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	arr := cloneViaYAML[[]any](cur[key])

	switch p.ItemOps {
	case "add":
		for _, v := range p.Value.([]any) {
			if !slices.Contains(arr, v) {
				arr = append(arr, v)
			}
		}

	case "remove":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(p.Value.([]any), v)
		})

	case "replace":
		for i, v := range arr {
			if v == p.Old {
				arr[i] = p.Value.([]any)[0]
				break
			}
		}

	case "keep":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return !slices.Contains(p.Value.([]any), v)
		})

	case "disable":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(p.Value.([]any), v)
		})
	}

	tgt[key] = arr
	return nil
}