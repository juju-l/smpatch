package smpatch//

import (
	"fmt"
	"strings"
	//
)

///** func

// Apply 对 obj 原地应用 patches
func Apply(obj map[string]any, patches []PatchOp) error {
	//
	for _, p := range patches {
	if err := applyOne(obj, p); err != nil {
	return err
	}
	}
	return nil
}
//

func applyOne(obj map[string]any, p PatchOp) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")
	parent, lastPart, err := resolveParent(obj, parts)
	if err != nil {
	return err
	}//
	switch p.Ope {
	case "merge":
	return applyMerge(parent, lastPart, p)
	case "replace":
	return applyRpl(parent, lastPart, p)
	case "delete":
	return applyDel(parent, lastPart, p)
	default:
	return fmt.Errorf("unknown ope: %s", p.Ope)
	//:
	//
	//
	}
	//
	//
	//
}

func resolveParent(obj map[string]any, parts []string) (map[string]any, string, error) {
	cur := obj
	for i, p := range parts[:len(parts)-1] {
	v, ok := cur[p]
	if !ok {
	return nil, "", fmt.Errorf("path not found: %s", strings.Join(parts[:i+1], "/")) ///**
	}
	m, ok := v.(map[string]any)
	if !ok {
	return nil, "", fmt.Errorf("intermediate node is not map")
	}
	cur = m
	}
	return cur, parts[len(parts)-1], nil
}

func init() {
	///**
}

// struct

// interface