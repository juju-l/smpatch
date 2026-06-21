package smpatch

import (
	"strings"
)

func mapAr(p *Patch, tgt map[string]any) error {
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

	for _, v := range p.Value.([]any) {
		m := v.(map[string]any)
		k := m[p.ByKey]
		found := false
		for i, e := range arr {
			if e.(map[string]any)[p.ByKey] == k {
				for mk, mv := range m {
					arr[i].(map[string]any)[mk] = mv
				}
				found = true
				break
			}
		}
		if !found {
			arr = append(arr, m)
		}
	}

	cur[key] = arr
	return nil
}