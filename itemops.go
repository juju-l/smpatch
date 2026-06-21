package smpatch

import "strings"

func itemOps(
	p *Patch,
	src map[string]any,
	tgt map[string]any,
) error {
	parts := strings.Split(strings.Trim(p.PathKey, "/"), "/")
	cur := src
	for i := 0; i < len(parts)-1; i++ {
		cur = cur[parts[i]].(map[string]any)
	}
	key := parts[len(parts)-1]

	arr := cur[key].([]any)
	vals := toSlice(p.Value)

	switch p.ItemOps {
	case "add":
		for _, v := range vals {
			found := false
			for _, e := range arr {
				if e == v {
					found = true
					break
				}
			}
			if !found {
				arr = append(arr, v)
			}
		}

	case "remove":
		newArr := []any{}
		for _, v := range arr {
			hit := false
			for _, r := range vals {
				if v == r {
					hit = true
					break
				}
			}
			if !hit {
				newArr = append(newArr, v)
			}
		}
		arr = newArr

	case "replace":
		for i, v := range arr {
			if v == p.Old {
				arr[i] = p.Value
				break
			}
		}

	case "keep":
		newArr := []any{}
		for _, v := range arr {
			for _, k := range vals {
				if v == k {
					newArr = append(newArr, v)
					break
				}
			}
		}
		arr = newArr

	case "disable":
		newArr := []any{}
		for _, v := range arr {
			hit := false
			for _, d := range vals {
				if v == d {
					hit = true
					break
				}
			}
			if !hit {
				newArr = append(newArr, v)
			}
		}
		arr = newArr
	}

	tgt[key] = arr
	return nil
}