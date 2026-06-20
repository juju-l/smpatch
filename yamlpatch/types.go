package main

//import "gopkg.in/yaml.v3"

// ResourcePatchSpec 对应 ResourcePatch CR 的 spec
type ResourcePatchSpec struct {
	Namespace string     `yaml:"namespace"`
	TargetRef TargetRef `yaml:"targetRef"`
	Patches  []PatchOp  `yaml:"patches"`
}

type TargetRef struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Name       string `yaml:"name"`
}

// PatchOp 单条 patch 操作
type PatchOp struct {
	Ope      string      `yaml:"ope"`
	PathKey  string      `yaml:"pathKey"`
	ByKey    string      `yaml:"byKey,omitempty"`
	ItemOps  string      `yaml:"itemOps,omitempty"`
	MixedAr  bool        `yaml:"mixedAr,omitempty"`
	Old      interface{} `yaml:"old,omitempty"`
	Value    interface{} `yaml:"value"`
}

// Diff 表示一次变更
type Diff struct {
	Path string
	Old  interface{}
	New  interface{}
}