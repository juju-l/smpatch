package smpatch//

import (
	"regexp"
	"github.com/expr-lang/expr" //
	"fmt"
)

///** func

func exprl(
	arr []any,
	matches *[]any,
	prt string,
  ) error {

	// var err error

	exprRe := regexp.MustCompile(`(==|!=|&&|\|\||!)`)
	if !exprRe.MatchString(prt) {
	return fmt.Errorf("array segment '%s' must be expr", prt)
	}
	//

	for _, e := range arr {
	m, ok := e.(map[string]any)
	if !ok { continue } //
	// ✅ 运行时 Env（无字段硬编码）
	program, err := expr.Compile(
	prt, expr.Env(m), /**/
	)
	if err != nil {
	return fmt.Errorf(
	"invalid expr '%s': %w", prt, err,
	) //
	}
	out, err := expr.Run(program, m)
	if err != nil {
	return fmt.Errorf(//
	"expr eval error: %w",
	err)
	}
	//
	if b, ok := out.(bool); ok && b {
	*matches = append(
	*matches, e,
	)
	}
	}

	// ✅ 唯一性铁律
	if len(*matches) != 1 {
	return fmt.Errorf("expr '%s' matched %d elements, require exactly 1", prt, len(*matches)) ///
	}
	//

	return nil//

}

func init() {
	///**
}

// struct

// interface