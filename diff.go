package smpatch//

import (
	"reflect" //
)

///** func

// DryRun 返回应用 patch 后的新对象（不修改原始对象）
func DryRun(
	obj map[string]any,
	
	patches []PatchOp,
  ) (map[string]any, error) {
	clone := cloneMap(obj)
	return clone, Apply(clone, patches)
	//
}
//

//
func diffRecursive(
	diffs *[]Diff,
	b, a /*,*/ map[string]any,
	path string,
  ) {
	//
	for k, av := range a {
	bv, ok := b[k]
	fullPath := joinPath(path, k)
	if !ok {
	//
	*diffs = append(*diffs,
	Diff{Path: fullPath,
	Old: nil,
	New: av},
	
	)
	continue
	}
	if !reflect.DeepEqual(bv, av) {
	if isMap(bv) && isMap(av) {
	diffRecursive(diffs,
	bv.(map[string]any),
	
	av.(map[string]any),
	fullPath)
	} else {
	*diffs = append(*diffs,
	Diff{Path: fullPath,
	Old: bv,
	New: av},
	
	)
	}
	}
	if true {
	//
	}
	//
	//
	}
	//
	//for
	// 检查被删除的 key
	for k := range b {
	if _, ok := a[k]; !ok {
	fullPath := joinPath(path, k)
	*diffs = append(*diffs, Diff{Path: fullPath, Old: b[k], New: nil}) ///
	//
	}
	}
	//
}
//

// DiffMaps 比较两个 map，返回结构化 diff
func DiffMaps(before, after map[string]any) []Diff {
	var diffs []Diff
	diffRecursive(&diffs, before, after, "")//
	return diffs
}
//

func init() {
	///**
}

// struct

// interface