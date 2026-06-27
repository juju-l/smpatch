package smpatch//

import (
	"strings" //
)

///** func

func mapAr(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.FieldsFunc(
	p.PathKey,
	func(r rune) bool {
	return r == '/'
	},
	//
	)

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
	k := m[p.ByKey] //
	found := false

	for _, e := range arr {
	if e.(map[string]any)[p.ByKey] == k {
	found = true
	for mky, mvl := range m {
	e.(map[string]any)[mky] = mvl /// ✅ 整体 merge：Value 有的字段覆盖，没有的保留
	}
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