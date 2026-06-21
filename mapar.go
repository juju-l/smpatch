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
		if cur[parts[i]] == nil {
			cur[parts[i]] = map[string]any{}
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	arr := cloneViaYAML[[]any](cur[key])

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

	tgt[key] = arr
	return nil
}