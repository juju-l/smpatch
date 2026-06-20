package main

import (
	"fmt"
)

func applyReplace(parent map[string]interface{}, key string, p PatchOp) error {
	if !p.MixedAr {
		return fmt.Errorf(
			"replace is only allowed for mixed arrays (mixedAr=true): path=%s",
			p.PathKey,
		)
	}

	// mixed 数组：直接覆盖
	parent[key] = p.Value
	return nil
}