package smpatch

import (
	"testing"
)

func TestApply_ScalarArray_Remove(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321, 987},
		},
	}

	tgt := map[string]any{}

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/p",
			ItemOps: "remove",
			Value:   []any{321},
		},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 1 || result[0] != 987 {
		t.Fatalf("remove failed, got %v", result)
	}
}