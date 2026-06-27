package smpatch//

import (
	"strings" //
)

///** func

func mixed(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.FieldsFunc(
	p.PathKey,func(r rune) bool { return r == '/' }, ////
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

	cur[key] = DeepCopy(p.Value)//

	return nil//

}

func init() {
	///**
}

// struct

// interface