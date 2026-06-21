package smpatch

import "strings"

func applyRpl(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			cur[parts[i]] = map[string]any{}
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	cur[key] = DeepCopy(p.Value)
	return nil
}