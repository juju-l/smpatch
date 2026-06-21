package smpatch

import (
	"testing"
)

func TestApply_ScalarArray_Remove(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "remove", Value: []any{321}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 1 || result[0] != 987 {
		t.Fatalf("remove failed, got %v", result)
	}
}

func TestApply_ScalarArray_Disable(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "disable", Value: []any{321}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 1 || result[0] != 987 {
		t.Fatalf("disable failed, got %v", result)
	}
}

func TestApply_ScalarArray_Add(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "add", Value: []any{987}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 2 {
		t.Fatalf("add failed, got %v", result)
	}
}

func TestApply_ScalarArray_Keep(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987, 123}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "keep", Value: []any{321, 987}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if len(result) != 2 {
		t.Fatalf("keep failed, got %v", result)
	}
}

func TestApply_ScalarArray_Replace(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"p": []any{321, 987}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/p", ItemOps: "replace", Old: 321, Value: []any{123}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["p"].([]any)
	if result[0] != 123 || result[1] != 987 {
		t.Fatalf("replace failed, got %v", result)
	}
}

func TestApply_StructArray_ByKey_Merge(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{"1", "2"}},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope: "merge", PathKey: "/spec/bindings", ByKey: "role",
			Value: []any{
				map[string]any{"role": "admin", "members": []any{"3"}},
			},
		},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	members := tgt["spec"].
		(map[string]any)["bindings"].
		([]any)[0].
		(map[string]any)["members"].([]any)

	if len(members) != 3 {
		t.Fatalf("byKey merge failed, got %v", members)
	}
}

func TestApply_MixedArray_Replace(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"tags": []any{1, "a", true}}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "replace", PathKey: "/spec/tags", MixedAr: true, Value: []any{"x", "y"}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	result := tgt["spec"].(map[string]any)["tags"].([]any)
	if len(result) != 2 || result[0] != "x" {
		t.Fatalf("mixed replace failed, got %v", result)
	}
}

func TestApply_NormalField_Merge(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"add": "r"}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec", Value: map[string]any{"c": "tst"}},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if tgt["spec"].(map[string]any)["c"] != "tst" {
		t.Fatalf("normal field merge failed")
	}
}

func TestApply_Delete_Field(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"add": "r"}}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "delete", PathKey: "/spec/add"},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	if _, ok := tgt["spec"].(map[string]any)["add"]; ok {
		t.Fatalf("delete field failed")
	}
}

func TestApply_Delete_FromStructArray(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "a"},
				map[string]any{"role": "b"},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{Ope: "delete", PathKey: "/spec/bindings", ByKey: "role", Value: "a"},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	bindings := tgt["spec"].(map[string]any)["bindings"].([]any)
	if len(bindings) != 1 {
		t.Fatalf("delete from struct array failed")
	}
}