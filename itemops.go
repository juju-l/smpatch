package smpatch

import (
	"fmt"
	"strings"
	"slices"

	"github.com/expr-lang/expr"
)

func itemOps(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	var walk func(cur any, idx int) error
	walk = func(cur any, idx int) error {
		// ✅ 最后一个 segment：执行 itemOps
		if idx == len(parts)-1 {
			if m, ok := cur.(map[string]any); ok {
				arr, ok := m[parts[idx]].([]any)
				if !ok {
					return fmt.Errorf("target field '%s' must be array", parts[idx])
				}

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

				m[parts[idx]] = arr
				return nil
			}
			return fmt.Errorf("unexpected type at final segment: %T", cur)
		}

		part := parts[idx]

		// ✅ map → 普通路径
		if m, ok := cur.(map[string]any); ok {
			next, ok := m[part]
			if !ok {
				return fmt.Errorf("path segment '%s' not found", part)
			}
			return walk(next, idx+1)
		}

		// ✅ array → 必须是表达式
		if arr, ok := cur.([]any); ok {
			if !isExpr(part) {
				return fmt.Errorf("array segment '%s' must be expr", part)
			}

			exprStr := normalizeExpr(part)

			program, err := expr.Compile(exprStr, expr.Env(map[string]any{
				"role": "",
				"env":  "",
			}))
			if err != nil {
				return fmt.Errorf("invalid expr '%s': %w", part, err)
			}

			var matches []any
			for _, e := range arr {
				m, ok := e.(map[string]any)
				if !ok {
					continue
				}
				out, err := expr.Run(program, m)
				if err != nil {
					return fmt.Errorf("expr eval error: %w", err)
				}
				if b, ok := out.(bool); ok && b {
					matches = append(matches, e)
				}
			}

			if len(matches) != 1 {
				return fmt.Errorf(
					"expr '%s' matched %d elements, require exactly 1",
					part, len(matches),
				)
			}

			return walk(matches[0], idx+1)
		}

		return fmt.Errorf("unexpected type at segment '%s': %T", part, cur)
	}

	return walk(tgt, 0)
}

func isExpr(s string) bool {
	return strings.Contains(s, "==") ||
		strings.Contains(s, "!=") ||
		strings.Contains(s, "&&") ||
		strings.Contains(s, "||") ||
		strings.HasPrefix(s, "!")
}

// ✅ 关键修复：递归处理 && / || 两边
func normalizeExpr(expr string) string {
	// 处理 ! 前缀
	if strings.HasPrefix(expr, "!") {
		inner := strings.TrimPrefix(expr, "!")
		return "!(" + normalizeExpr(inner) + ")"
	}

	// 处理 &&
	if strings.Contains(expr, "&&") {
		parts := strings.SplitN(expr, "&&", 2)
		return normalizeExpr(parts[0]) + " && " + normalizeExpr(parts[1])
	}

	// 处理 ||
	if strings.Contains(expr, "||") {
		parts := strings.SplitN(expr, "||", 2)
		return normalizeExpr(parts[0]) + " || " + normalizeExpr(parts[1])
	}

	// 处理 == / !=
	for _, op := range []string{"==", "!="} {
		parts := strings.SplitN(expr, op, 2)
		if len(parts) == 2 {
			left := strings.TrimSpace(parts[0])
			right := strings.TrimSpace(parts[1])
			if right != "" && right[0] != '"' && right[0] != '\'' {
				return left + " " + op + " \"" + right + "\""
			}
		}
	}

	return expr
}