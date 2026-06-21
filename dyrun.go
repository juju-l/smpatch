package smpatch

func DyRun(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
) error {
	for k, v := range src {
		tgt[k] = cloneViaYAML[[]any](v)
	}
	return Apply(src, patches, tgt)
}