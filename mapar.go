package smpatch//

import (
	"strings"
	. "slices"
	//
)

///** func

func mapAr(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.Split(
	strings.Trim(p.PathKey, "/"),//
	"/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
	//
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

	for _, e := range arr {
	if e.(map[string]any)[p.ByKey] == k {
	//
	// ✅ 只 merge，不替换
	for mk, mv := range m {
	if mk == "members" {
	// ✅ 数组合并
	oldMembers := e.(map[string]any)["members"].([]any) ///
	newMembers := mv.([]any)
	for _, nm := range newMembers {
	if !Contains(oldMembers, nm) {
	oldMembers = append(oldMembers, nm)
	}
	}
	e.(map[string]any)["members"] = oldMembers
	} else {
	e.(map[string]any)[mk] = mv
	}
	}
	found = true
	break
	}
	}

	if !found {
	arr = append(arr, m)//
	}
	}

	cur[key] = arr

	return nil//

}

func init() {
	///**
}

// struct

// interface