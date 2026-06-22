package smpatch

import (
	"fmt"
	"strings"
	"slices"
)

func itemOps(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	cur := tgt
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]

		// ✅ 当前是数组，且下一段是表达式 → 筛选
		if arr, ok := cur[part].([]any); ok {
			if i+1 >= len(parts) {
				return fmt.Errorf("unexpected end of path after array")
			}

			exprPart := parts[i+1]
			if !strings.ContainsAny(exprPart, "=&|!") {
				return fmt.Errorf("array segment '%s' must be followed by expr", part)
			}

			var matches []any
			for _, e := range arr {
				m, ok := e.(map[string]any)
				if !ok {
					continue
				}
				if evalSimpleExpr(m, exprPart) {
					matches = append(matches, e)
				}
			}

			if len(matches) != 1 {
				return fmt.Errorf(
					"expr '%s' matched %d elements, require exactly 1",
					exprPart, len(matches),
				)
			}

			cur = matches[0].(map[string]any)
			i++ // ✅ 跳过已处理的表达式段
			continue
		}

		// ✅ 普通 struct 路径
		next, ok := cur[part].(map[string]any)
		if !ok {
			return fmt.Errorf("path segment '%s' is not a struct", part)
		}
		cur = next
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
	if strings.HasPrefix(expr, "!") {
		return !evalSimpleExpr(obj, expr[1:])
	}

	if strings.Contains(expr, "&&") {
		parts := strings.SplitN(expr, "&&", 2)
		return evalSimpleExpr(obj, parts[0]) &&
			evalSimpleExpr(obj, parts[1])
	}

	if strings.Contains(expr, "||") {
		parts := strings.SplitN(expr, "||", 2)
		return evalSimpleExpr(obj, parts[0]) ||
			evalSimpleExpr(obj, parts[1])
	}

	// ✅ 正确解析 == / !=
	eq := strings.Index(expr, "==")
	ne := strings.Index(expr, "!=")

	var key, op, want string
	if eq != -1 {
		key = expr[:eq]
		op = "=="
		want = expr[eq+2:]
	} else if ne != -1 {
		key = expr[:ne]
		op = "!="
		want = expr[ne+2:]
	} else {
		return false
	}

	actual := fmt.Sprint(obj[key])
	switch op {
	case "==":
		return actual == want
	case "!=":
		return actual != want
	}
	return false
}

// struct

// interface