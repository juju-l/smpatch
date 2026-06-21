package smpatch//

import (
	"fmt"//
	"strings"
	"slices"
)

///** func

func itemOps(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// ✅ 表达式路径段
		if strings.Contains(part, "==") ||
			strings.Contains(part, "&&") ||
			strings.Contains(part, "||") ||
			strings.HasPrefix(part, "!") {

			arr, ok := cur[part].([]any)
			if !ok {
				return fmt.Errorf("expr target is not array")
			}

			var matches []any
			for _, e := range arr {
				m, ok := e.(map[string]any)
				if !ok {
					continue
				}
				if evalSimpleExpr(m, part) {
					matches = append(matches, e)
				}
			}

			if len(matches) != 1 {
				return fmt.Errorf(
					"expr '%s' matched %d elements, require exactly 1",
					part, len(matches),
				)
			}

			cur = matches[0].(map[string]any)
			continue
		}

		// ✅ 普通路径
		if cur[part] == nil {
			cur[part] = map[string]any{}
		}
		cur = cur[part].(map[string]any)
	}

	key := parts[len(parts)-1]

	arr := DeepCopy(cur[key]).([]any)
	vals := p.Value.([]any)

	switch p.ItemOps {
	case "add":
		for _, v := range vals {
			if !slices.Contains(arr, v) {
				arr = append(arr, v)
			}
		}
	case "remove":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(vals, v)
		})
	case "replace":
		for i, v := range arr {
			if v == p.Old {
				arr[i] = vals[0]
				break
			}
		}
	case "keep":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return !slices.Contains(vals, v)
		})
	case "disable":
		arr = slices.DeleteFunc(arr, func(v any) bool {
			return slices.Contains(vals, v)
		})
	}

	cur[key] = arr
	return nil
}

func init() {
	///**
}

func evalSimpleExpr(obj map[string]any, expr string) bool {
	// !role==admin
	if strings.HasPrefix(expr, "!") {
		return !evalSimpleExpr(obj, expr[1:])
	}

	// role==admin && env==prod
	if strings.Contains(expr, "&&") {
		parts := strings.SplitN(expr, "&&", 2)
		return evalSimpleExpr(obj, parts[0]) &&
			evalSimpleExpr(obj, parts[1])
	}

	// role==admin || role==editor
	if strings.Contains(expr, "||") {
		parts := strings.SplitN(expr, "||", 2)
		return evalSimpleExpr(obj, parts[0]) ||
			evalSimpleExpr(obj, parts[1])
	}

	// role==admin / role!=admin
	kv := strings.SplitN(expr, "=", 2)
	if len(kv) != 2 {
		return false
	}

	actual := fmt.Sprint(obj[kv[0]])
	want := kv[1]

	if kv[0] == "!" {
		return actual != want
	}
	return actual == want
}

// struct

// interface