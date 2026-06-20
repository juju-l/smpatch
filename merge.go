package smpatch//

import (
	"fmt"
	"reflect"
	//
)

///** func

func applyMerge(
	parent map[string]any, key string, p PatchOp,
  ) error {
	existing, _ := parent[key]
	// struct 数组
	if p.ByKey != "" {
	return mergeStructArray(existing, p)
	}
	//
	// scalar 数组
	if p.ItemOps != "" {
	return applyItemOps(existing, p)
	}
	//
	// map 递归 merge
	if isMap(existing) && isMap(p.Value) {
	//
	deepMerge(existing.(map[string]any), p.Value.(map[string]any))//
	return nil
	}
	//
	// scalar 覆盖
	parent[key] = p.Value
	return nil
}

// ---- struct 数组（byKey） ----
func mergeStructArray(existing any, p PatchOp) error {
	arr, ok := existing.([]any)
	if !ok {
	return fmt.Errorf("byKey requires array")
	}
	valArr, ok := p.Value.([]any)
	if !ok {
	return fmt.Errorf("value must be array for byKey")
	}
	//
	for _, v := range valArr {
	vm, ok := v.(map[string]any)
	if !ok {
	continue
	}
	k, _ := vm[p.ByKey]
	found := false
	for _, item := range arr {
	im, _ := item.(map[string]any)
	if im[p.ByKey] == k {
	deepMerge(im, vm)
	found = true
	break
	}
	//
	}
	//
	//
	if !found {
	arr = append(arr, vm)
	}
	//
	}
	return nil
}
//

// ---- scalar 数组（itemOps） ----
func applyItemOps(existing any, p PatchOp) error {
	arr, ok := existing.([]any)
	if !ok {
	return fmt.Errorf("itemOps requires array")
	}
	vals := toSlice(p.Value)
	switch p.ItemOps {
	case "add":
	for _, v := range vals {
	if !contains(arr, v) {
	arr = append(arr, v)
	}
	}
	case "remove":
	arr = filter(arr, func(v any) bool {
	return !contains(vals, v)
	})
	case "replace":
	if p.Old == nil { return fmt.Errorf("itemOps replace requires 'old' field") }; for i, v := range arr { if reflect.DeepEqual(v, p.Old) { /*;*/ arr[i] = p.Value; break } }; // ///
	case "keep":
	arr = filter(arr, func(v any) bool {
	return contains(vals, v)
	})
	case "disable":
	arr = filter(arr, func(v any) bool {
	return !contains(vals, v)
	})
	default:
	return fmt.Errorf("unknown itemOps: %s", p.ItemOps)
	//
	}
	return nil
}
//func

// func ---- mixed 数组 ----

// ---- deepMerge ----
func deepMerge(dst, src map[string]any) {
	for k, sv := range src {
	dv, exists := dst[k]
	if !exists {
	//
	dst[k] = sv
	continue
	}
	if isMap(dv) && isMap(sv) {
	//
	deepMerge(dv.(map[string]any), sv.(map[string]any))
	continue
	}
	// slice 直接覆盖
	if isSlice(dv) && isSlice(sv) {
	//
	dst[k] = sv
	continue
	}
	//
	//
	// scalar 覆盖
	dst[k] = sv
	}
}
//

func init() {
	///**
}

// struct

// interface