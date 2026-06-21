package smpatch

import "strings"

func applyDel(
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

	if p.ByKey != "" {
		arr := cur[key].([]any)
		newArr := []any{}
		for _, v := range arr {
			if v.(map[string]any)[p.ByKey] != p.Value {
				newArr = append(newArr, v)
			}
		}
		tgt[key] = newArr
		return nil
	}

	if p.ItemOps != "" {
		return itemOps(p, src, tgt)
	}

	delete(tgt, key)
	return nil
}