package smpatch

func DiffCmp(tgt map[string]any, patches []*Patch) error {
	return DyRun(tgt, patches)
}