package smpatch

import (
	"strings"
	"slices"
)

func applyDel(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			return nil
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	if p.ByKey != "" {
		arr := cur[key].([]any)
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return v.(map[string]any)[p.ByKey] == p.Value
		})
		cur[key] = arr
		return nil
	}

	if p.ItemOps != "" {
		return itemOps(p, tgt)
	}

	delete(cur, key)
	return nil
}