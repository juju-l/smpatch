package smpatch

import (
	"testing"
)

func clone(src any) any {
	switch v := src.(type) {

	case map[string]any:
		dst := map[string]any{}
		for k, val := range v {
			dst[k] = clone(val)
		}
		return dst

	case []any:
		dst := make([]any, len(v))
		for i, val := range v {
			dst[i] = clone(val)
		}
		return dst

	default:
		return v
	}
}

// ============================================
// Scalar Array Tests (/spec/p)
// ============================================

func TestApply_ScalarArray_Remove(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "remove", Value: []any{321}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 1 || result[0] != 987 {
		t.Fatalf("remove failed, got %v", result)
	}
}

func TestApply_ScalarArray_Disable(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "disable", Value: []any{321}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 1 || result[0] != 987 {
		t.Fatalf("disable failed, got %v", result)
	}
}

func TestApply_ScalarArray_Add(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321}}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "add", Value: []any{987}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 2 {
		t.Fatalf("add failed, got %v", result)
	}
}

func TestApply_ScalarArray_Keep(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987, 123}}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "keep", Value: []any{321, 987}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 2 {
		t.Fatalf("keep failed, got %v", result)
	}
}

func TestApply_ScalarArray_Replace(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "replace", Old: 321, Value: 123},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if result[0] != 123 || result[1] != 987 {
		t.Fatalf("replace failed, got %v", result)
	}
}

// ============================================
// Struct Array Tests (/spec/bindings)
// ============================================

func TestApply_StructArray_ByKey_Merge(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "qwertyuiophjkl", "members": []any{"1", "2"}},
			},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge", PathKey: "/spec/bindings", ByKey: "role",
			Value: []any{
				map[string]any{"role": "qwertyuiophjkl", "members": []any{"3", "yuiop"}},
			},
		},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	bindings := tgt["spec"].(map[string]any)["bindings"].([]any)
	members := bindings[0].(map[string]any)["members"].([]any)

	// 验证合并：原有 + 新增
	if len(members) != 4 {
		t.Fatalf("byKey merge failed, got %v", members)
	}
}

// ============================================
// Mixed Array & Normal Field Tests
// ============================================

func TestApply_MixedArray_Replace(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"tags": []any{1, "a", true}}}
	patches := []*Patch{
		{Ope: "replace", PathKey: "/spec/tags", MixedAr: true, Value: []any{"x", "y"}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	result := tgt["spec"].(map[string]any)["tags"].([]any)
	if len(result) != 2 || result[0] != "x" {
		t.Fatalf("mixed replace failed, got %v", result)
	}
}

func TestApply_NormalField_Merge(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"add": "r"}}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec", Value: map[string]any{"c": "tst"}},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	if tgt["spec"].(map[string]any)["c"] != "tst" {
		t.Fatalf("normal field merge failed")
	}
}

// ============================================
// Delete Tests
// ============================================

func TestApply_Delete_Field(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"add": "r"}}
	patches := []*Patch{
		{Ope: "delete", PathKey: "/spec/add"},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	if _, ok := tgt["spec"].(map[string]any)["add"]; ok {
		t.Fatalf("delete field failed")
	}
}

func TestApply_Delete_FromStructArray(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "qwertyuiopzsd"},
				map[string]any{"role": "qwertyuiophjkl"},
			},
		},
	}
	patches := []*Patch{
		{Ope: "delete", PathKey: "/spec/bindings", ByKey: "role", Value: "qwertyuiopzsd"},
	}
	tgt := clone(src).(map[string]any)
	Apply(src, patches, tgt)

	bindings := tgt["spec"].(map[string]any)["bindings"].([]any)
	if len(bindings) != 1 {
		t.Fatalf("delete from struct array failed")
	}
}
