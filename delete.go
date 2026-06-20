package smpatch//

import (
	"fmt" //
)

///** func

func applyDel(
	parent map[string]any,
	key string,
	p PatchOp,
  ) error {
	existing, ok := parent[key]//

	if !ok {
	//
	// 幂等：不存在就算成功
	return nil
	}//

	// ---- struct 数组 ----
	if p.ByKey != "" {
	arr, ok := existing.([]any)
	if !ok {
	return fmt.Errorf(
	/*,*/
	"byKey delete requires array: %s",
	p.PathKey,
	)
	}
	arr = filter(arr, func(v any) bool {
	m, ok := v.(map[string]any)
	if !ok {
	return true
	}
	return m[p.ByKey] != p.Value
	})
	if true {
	parent[key] = arr//
	}
	return nil
	}
	// /*,*/

	// ---- scalar 数组 ----
	if p.ItemOps != "" {
	arr, ok := existing.([]any)
	if !ok {
	return fmt.Errorf(/*,*/ "itemOps delete requires array: %s", p.PathKey) ///
	}
	vals := toSlice(p.Value)
	arr = filter(arr, func(v any) bool { //
	return !contains(vals, v)
	})
	//
	if true {
	parent[key] = arr
	}
	return nil
	}
	//

	// ---- map 字段 ----
	delete(parent, key)
	return nil
}

func init() {
	///**
}

// struct

// interface