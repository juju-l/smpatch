package smpatch

import "strings"

func dpMeg(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		if cur[parts[i]] == nil {
			cur[parts[i]] = map[string]any{}
		}
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	switch v := cur[key].(type) {
	case map[string]any:
		cp := DeepCopy(v).(map[string]any)
		for mk, mv := range p.Value.(map[string]any) {
			cp[mk] = mv
		}
		cur[key] = cp
	default:
		cur[key] = DeepCopy(p.Value)
	}
	return nil
}