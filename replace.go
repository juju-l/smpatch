package smpatch//

import (
	"fmt" //
)

///** func

func applyRpl(
	parent map[string]any,
	key string,
	p PatchOp,
  ) error {
	//
	//
	//replace
	//
	//
	if !p.MixedAr {
	return fmt.Errorf(
	/*,*/
	"replace is only allowed for mixed arrays (mixedAr=true): path=%s", ///
	p.PathKey,
	)
	}
	// mixed 数组：直接覆盖
	parent[key] = p.Value
	return nil
}

func init() {
	///**
}

// struct

// interface