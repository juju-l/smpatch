package main

import (
	"fmt"
)

func applyDelete(parent map[string]interface{}, key string, p PatchOp) error {
	existing, ok := parent[key]
	if !ok {
		// 幂等：不存在就算成功
		return nil
	}

	// ---- struct 数组 ----
	if p.ByKey != "" {
		arr, ok := existing.([]interface{})
		if !ok {
			return fmt.Errorf("byKey delete requires array: %s", p.PathKey)
		}
		arr = filter(arr, func(v interface{}) bool {
			m, ok := v.(map[string]interface{})
			if !ok {
				return true
			}
			return m[p.ByKey] != p.Value
		})
		parent[key] = arr
		return nil
	}

	// ---- scalar 数组 ----
	if p.ItemOps != "" {
		arr, ok := existing.([]interface{})
		if !ok {
			return fmt.Errorf("itemOps delete requires array: %s", p.PathKey)
		}
		vals := toSlice(p.Value)
		arr = filter(arr, func(v interface{}) bool {
			return !contains(vals, v)
		})
		parent[key] = arr
		return nil
	}

	// ---- map 字段 ----
	delete(parent, key)
	return nil
}