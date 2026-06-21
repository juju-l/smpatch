package smpatch

import (
	"fmt"
	//"bytes"
	"gopkg.in/yaml.v3"
)

type Patch struct {
	Ope     string `yaml:"ope"`
	PathKey string `yaml:"pathKey"`

	ByKey   string `yaml:"byKey,omitempty"`
	ItemOps string `yaml:"itemOps,omitempty"`
	MixedAr bool   `yaml:"mixedAr,omitempty"`

	Old   any `yaml:"old,omitempty"`
	Value any `yaml:"value"`
}

// cloneViaYAML 使用 YAML Marshal / Unmarshal 实现深拷贝
// 泛型 T 仅用于约束返回类型，内部实现使用 any
func cloneViaYAML[T any](v any) T {
	data, err := yaml.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("yaml marshal failed: %v", err))
	}

	var out T
	if err := yaml.Unmarshal(data, &out); err != nil {
		panic(fmt.Sprintf("yaml unmarshal failed: %v", err))
	}

	return out
}