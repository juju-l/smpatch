package smpatch

func DyRun(tgt map[string]any, patches []*Patch) error {
	return Apply(tgt, patches, tgt)
}