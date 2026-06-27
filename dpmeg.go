package smpatch//

import (
	"strings" //
)

///** func

func dpMeg(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.FieldsFunc(
	p.PathKey,func(
	r rune,
	) bool { return r == '/' }, ////
	)

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
	//
	if cur[parts[i]] == nil {
	cur[parts[i]] = map[string]any{}
	}
	cur = 
	cur[parts[i]].
	(map[string]any)
	}
	key := parts[len(parts)-1]

	switch v := cur[key].(type) {
	// case：
			//
	  //
	case map[string]any:
	cpy := DeepCopy(v).(map[string]any)
	for mky, mvl := range p.Value.(map[string]any) {
	cpy[mky] = DeepCopy(mvl)
	}
	cur[key] = cpy
	  //
	default:
	cur[key] = DeepCopy(p.Value)
	  //
	}

	return nil//

}

func init() {
	///**
}

// struct

// interface