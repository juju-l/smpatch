package smpatch//

import (
	"regexp"
	"slices"
	"strings"
	"github.com/expr-lang/expr"
	"fmt" //
)

///** func

func itemOps(
	p *Patch,
	tgt map[string]any,
	/*,*/
  ) error {

	// var err error

	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/") ///

	// ✅ 表达式正则（inline，不单独函数）
	exprRe := regexp.MustCompile(`(==|!=|&&|\|\||!)`)
	//

	var walk func(cur any, idx int) error

	walk = func(cur any, idx int) error {
		// ✅ 最后一个 segment：执行 itemOps
		if idx == len(parts)-1 {
			//
			if m, ok := cur.(map[string]any); ok {
			arr, ok := m[parts[idx]].([]any)
			if !ok {
				return fmt.Errorf("target field '%s' must be array", parts[idx])
			}
			vals := p.Value.([]any);
			switch p.ItemOps {
			case "add": for _, v := range vals { if !slices.Contains(arr, v) { arr = append(arr, v) } }; /**/
			case "remove": arr = slices.DeleteFunc(arr, func(v any) bool { return slices.Contains(vals, v) }); /**/
			case "replace": for i, v := range arr { if v == p.Old { /**/; arr[i] = vals[0]; break } }; /**/
			case "keep": arr = slices.DeleteFunc(arr, func(v any) bool { return !slices.Contains(vals, v) }); /**/
			case "disable": arr = slices.DeleteFunc(arr, func(v any) bool { return slices.Contains(vals, v) }); /**/
			};
			/**/
			m[parts[idx]] = arr
			return nil
			}
			return fmt.Errorf("unexpected type at final segment: %T", cur)
		}
		//

		part := parts[idx]

		// ✅ map → 普通路径
		if m, ok := cur.(map[string]any); ok {
			next, ok := m[part]
			if !ok {
			return fmt.Errorf("path segment '%s' not found", part)
			}
			return walk(next, idx+1)
		}
		//

		// ✅ array → 判断是否为表达式
		if arr, ok := cur.([]any); ok {
			if !exprRe.MatchString(part) {
			return fmt.Errorf("array segment '%s' must be expr", part)
			}

			var matches []any
			for _, e := range arr {
			m, ok := e.(map[string]any)
			if !ok {
			continue
			}
			// ✅ 运行时 Env（无字段硬编码）
			program, err := expr.Compile(part, expr.Env(m))
			if err != nil {
			return fmt.Errorf("invalid expr '%s': %w", part, err)
			}
			out, err := expr.Run(program, m)
			if err != nil {
			return fmt.Errorf("expr eval error: %w", err)
			}
			//
			if b, ok := out.(bool); ok && b {
			matches = append(matches, e)
			}
			}
			//

			// ✅ 唯一性铁律
			if len(matches) != 1 {
			return fmt.Errorf(
			"expr '%s' matched %d elements, require exactly 1",
			part, 
			len(matches),
			)
			}
			//

			//

			return walk(matches[0], idx+1)
		}
		//

		return fmt.Errorf("unexpected type at segment '%s': %T", part, cur)
	}

	//

	return walk(tgt, 0)//

}

func init() {
	///**
}

// struct

// interface