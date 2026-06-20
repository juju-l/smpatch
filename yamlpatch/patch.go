package main

import (
	"fmt"
	"strings"
)

// Apply 对 obj 原地应用 patches
func Apply(obj map[string]interface{}, patches []PatchOp) error {
	for _, p := range patches {
		if err := applyOne(obj, p); err != nil {
			return err
		}
	}
	return nil
}

func applyOne(obj map[string]interface{}, p PatchOp) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")
	parent, lastPart, err := resolveParent(obj, parts)
	if err != nil {
		return err
	}

	switch p.Ope {
	case "merge":
		return applyMerge(parent, lastPart, p)
	case "replace":
		return applyReplace(parent, lastPart, p)
	case "delete":
		return applyDelete(parent, lastPart, p)
	default:
		return fmt.Errorf("unknown ope: %s", p.Ope)
	}
}

func resolveParent(obj map[string]interface{}, parts []string) (map[string]interface{}, string, error) {
	cur := obj
	for i, p := range parts[:len(parts)-1] {
		v, ok := cur[p]
		if !ok {
			return nil, "", fmt.Errorf("path not found: %s", strings.Join(parts[:i+1], "/"))
		}
		m, ok := v.(map[string]interface{})
		if !ok {
			return nil, "", fmt.Errorf("intermediate node is not map")
		}
		cur = m
	}
	return cur, parts[len(parts)-1], nil
}