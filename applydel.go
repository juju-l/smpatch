package smpatch

import (
	"slices"
	"strings"
	//"slices"
)

func applyDel(
	p *Patch,
	src map[string]any,
	tgt map[string]any,
) error {

	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			return nil // ✅ 不存在就是成功
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
		return itemOps(p, src, tgt)
	}

	delete(cur, key)
	return nil
}