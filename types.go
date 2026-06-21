package smpatch

type Patch struct {
	Ope     string `yaml:"ope"`
	PathKey string `yaml:"pathKey"`

	ByKey   string `yaml:"byKey,omitempty"`
	ItemOps string `yaml:"itemOps,omitempty"`
	MixedAr bool   `yaml:"mixedAr,omitempty"`

	Old   any `yaml:"old,omitempty"`
	Value any `yaml:"value"`
}