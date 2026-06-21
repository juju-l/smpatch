package smpatch

import (
	"strings"
	"slices"
)

func applyDel(
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

	if p.ByKey != "" {
		arr := cloneViaYAML[[]any](cur[key])
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return v.(map[string]any)[p.ByKey] == p.Value
		})
		tgt[key] = arr
		return nil
	}

	if p.ItemOps != "" {
		return itemOps(p, src, tgt)
	}

	delete(tgt, key)
	return nil
}