package smpatch

import "strings"

func mapAr(
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

	arr := cur[key].([]any)
	for _, v := range p.Value.([]any) {
		m := v.(map[string]any)
		k := m[p.ByKey]
		for i, e := range arr {
			if e.(map[string]any)[p.ByKey] == k {
				for mk, mv := range m {
					e.(map[string]any)[mk] = mv
				}
				goto NEXT
			}
		}
		arr = append(arr, m)
	NEXT:
	}

	tgt[key] = arr
	return nil
}