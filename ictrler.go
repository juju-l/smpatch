package smpatch//

import (
	"slices"
	// "strings"
	"fmt"
)

///** func

func ictrler(
	p *Patch,
	cur     any,
	s   string,
  ) error {

	// var err error

	if m, ok := cur.(map[string]any); ok {
	arr, ok := m[s].([]any)
	if !ok {
		return fmt.Errorf("target field '%s' must be array", s)
	}
	vals := p.Value.([]any);
	switch p.ItemOps {
	case "add": for _, v := range vals { if !slices.Contains(arr, v) { arr = append(arr, v) } }; /**/
	case "remove": arr = slices.DeleteFunc(arr, func(v any) bool { return slices.Contains(vals, v) }); /**/
	case "replace": for i, v := range arr { if v == p.Old { /**/; arr[i] = vals[0]; break } }; /**/
	case "disable": arr = slices.DeleteFunc(arr, func(v any) bool { return slices.Contains(vals, v) }); /**/ ///
	case "keep": arr = slices.DeleteFunc(arr, func(v any) bool { return !slices.Contains(vals, v) }); /**/
	};
	/**/
	m[s] = arr
	return nil
	}

	return fmt.Errorf("unexpected type at final segment: %T", cur)//

}

func init() {
	///**
}

// struct

// interface