package smpatch

import (
	"fmt"
	"strings"
	"slices"

	"github.com/expr-lang/expr"
)

func itemOps(p *Patch, tgt map[string]any) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")

	// ✅ inline 递归 walk
	var walk func(cur any, idx int) error
	walk = func(cur any, idx int) error {
		// 到达最后一个 segment → 执行 itemOps
		if idx == len(parts)-1 {
			arr := DeepCopy(cur).([]any)
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
			// 回溯：把结果写回父级
			// 这里 cur 已经是目标数组，直接返回
			return nil
		}

		part := parts[idx]

		// 当前层级一定是 struct 数组（大前提）
		arr, ok := cur.([]any)
		if !ok {
			return fmt.Errorf("expected array at segment '%s', got %T", part, cur)
		}

		// 判断是否为表达式 segment
		if isExpr(part) {
			// 编译表达式
			program, err := expr.Compile(part, expr.Env(map[string]any{
				"role": "",
				"env":  "",
			}))
			if err != nil {
				return fmt.Errorf("invalid expr '%s': %w", part, err)
			}

			// 筛选
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

			// 唯一性校验
			if len(matches) != 1 {
				return fmt.Errorf(
					"expr '%s' matched %d elements, require exactly 1",
					part, len(matches),
				)
			}

			// 唯一匹配 → 继续 walk 下一个 segment
			return walk(matches[0], idx+1)
		}

		// 非表达式：按普通路径处理
		// 找到名为 part 的字段（在 struct 数组中，这意味着 part 是字段名，需要是单个 struct）
		// 但按 PathKey 语义，非表达式 segment 应该是 struct 的字段名
		// 这里 cur 是数组，part 是字段名 → 需要先在数组中定位，再取字段
		// 按你的设计：非表达式 segment 不会出现在数组层，所以这里应该是 struct
		return fmt.Errorf("unexpected non-expr segment '%s' in array context", part)
	}

	// 启动递归
	return walk(tgt, 0)
}

// isExpr 判断 segment 是否为表达式
func isExpr(segment string) bool {
	return strings.Contains(segment, "==") ||
		strings.Contains(segment, "!=") ||
		strings.Contains(segment, "&&") ||
		strings.Contains(segment, "||") ||
		strings.HasPrefix(segment, "!")
}