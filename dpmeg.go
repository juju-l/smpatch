package smpatch

import "strings"

func dpMeg(
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

	switch v := cur[key].(type) {
	case map[string]any:
		cp := cloneViaYAML[map[string]any](v)
		for mk, mv := range p.Value.(map[string]any) {
			cp[mk] = mv
		}
		tgt[key] = cp
	default:
		tgt[key] = cloneViaYAML[any](p.Value)
	}

	return nil
}