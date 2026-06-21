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
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	switch v := cur[key].(type) {
	case map[string]any:
		for mk, mv := range p.Value.(map[string]any) {
			v[mk] = mv
		}
		tgt[key] = v
	default:
		tgt[key] = p.Value
	}

	return nil
}