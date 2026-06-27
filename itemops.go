package smpatch//

import (
	"strings"
	//
	"fmt"
)

///** func

func itemOps(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.Split(
	strings.Trim(p.PathKey, "/"),//
	"/")

	var walk func(cur any, idx int) error
	//
	walk = func(cur any, idx int) error {
	//
	// ✅ 最后一个 segment：执行 itemOps
	if idx == len(parts)-1 {
	return ictrler(p, cur, parts[idx])
	}
	//
	part := parts[idx]
	// ✅ map → 普通路径
	if m, ok := 
	cur.(map[string]any); 
	ok {
	next, ok := m[part]
	if !ok {
	return fmt.Errorf(//
	"path segment '%s' not found",
	part)
	}
	return walk(next, idx+1)
	}
	//
	// ✅ array → 判断是否为表达式
	if arr, ok := cur.([]any); ok {
	var matches []any
	if err := exprl(
	arr, &matches, part,
	); err != nil {
	return err
	}
	return walk(matches[0], idx+1)
	}
	//
	return fmt.Errorf("unexpected type at segment '%s': %T", part, cur) ///
	//
	}
	//
	return walk(tgt, 0)//
	
	//
	
	//

}

func init() {
	///**
}

// struct

// interface