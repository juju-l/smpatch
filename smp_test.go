package smpatch//

import (
	"testing" //
)

///** func

// ==================== 表达式 PathKey 测试 ====================

func TestItemOps_Expr_Equal(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{1, 2, 3}},
				map[string]any{"role": "viewer", "members": []any{4, 5}},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role==admin/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 2 {
		t.Fatalf("remove failed, got %v", members)
	}
}

// func TestItemOps_Expr_NotEqual(t *testing.T) {
// 	src := map[string]any{
// 		"spec": map[string]any{
// 			"bindings": []any{
// 				map[string]any{"role": "admin", "members": []any{1, 2}},
// 				map[string]any{"role": "viewer", "members": []any{3, 4}},
// 			},
// 		},
// 	}
// 	tgt := DeepCopy(src).(map[string]any)

// 	patches := []*Patch{
// 		{
// 			Ope:     "merge",
// 			PathKey: "/spec/bindings/role!=viewer/members",
// 			ItemOps: "remove",
// 			Value:   []any{1},
// 		},
// 	}
// 	if err := Apply(src, patches, tgt); err != nil {
// 		t.Fatalf("Apply failed: %v", err)
// 	}
// 	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
// 	if len(members) != 1 {
// 		t.Fatalf("remove with != failed, got %v", members)
// 	}
// }

// func TestItemOps_Expr_And(t *testing.T) {
// 	src := map[string]any{
// 		"spec": map[string]any{
// 			"bindings": []any{
// 				map[string]any{"role": "admin", "env": "prod", "members": []any{1, 2, 3}},
// 				map[string]any{"role": "admin", "env": "dev", "members": []any{4, 5}},
// 			},
// 		},
// 	}
// 	tgt := DeepCopy(src).(map[string]any)

// 	patches := []*Patch{
// 		{
// 			Ope:     "merge",
// 			PathKey: "/spec/bindings/role==admin && env==prod/members",
// 			ItemOps: "remove",
// 			Value:   []any{1},
// 		},
// 	}
// 	if err := Apply(src, patches, tgt); err != nil {
// 		t.Fatalf("Apply failed: %v", err)
// 	}
// 	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
// 	if len(members) != 2 {
// 		t.Fatalf("remove with && failed, got %v", members)
// 	}
// }

func TestItemOps_Expr_Or(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{1, 2}},
				map[string]any{"role": "viewer", "members": []any{3, 4}},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role==admin || role==editor/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("remove with || failed, got %v", members)
	}
}

func TestItemOps_Expr_NotPrefix(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{1, 2}},
				map[string]any{"role": "viewer", "members": []any{3, 4}},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/!role==viewer/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("remove with ! prefix failed, got %v", members)
	}
}

func TestItemOps_Expr_NoMatch(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{1, 2}},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role==nonexistent/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for no match")
	}
}

// func TestItemOps_Expr_NotUnique(t *testing.T) {
// 	src := map[string]any{
// 		"spec": map[string]any{
// 			"bindings": []any{
// 				map[string]any{"role": "admin"},
// 				map[string]any{"role": "admin"},
// 			},
// 		},
// 	}
// 	tgt := DeepCopy(src).(map[string]any)

// 	patches := []*Patch{
// 		{
// 			Ope:     "merge",
// 			PathKey: "/spec/bindings/role==admin/members",
// 			ItemOps: "remove",
// 			Value:   []any{1},
// 		},
// 	}
// 	if err := Apply(src, patches, tgt); err == nil {
// 		t.Fatal("expected error for non-unique expr result")
// 	}
// }

func TestItemOps_Expr_InvalidSyntax(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin"},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role===admin/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for invalid expr syntax")
	}
}

func TestApply_NormalField_Merge(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"add": "r"},
	}
	patches := []*Patch{
			{Ope: "merge", PathKey: "/spec",
			Value: map[string]any{"c": "tst"},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil { //
		t.Fatalf("Apply failed: %v", err)
	}
	if tgt["spec"].(map[string]any)["c"] != "tst" {
		t.Fatalf("normal field merge failed")
	}
}

func TestApply_StructArray_ByKey_Merge(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"bindings": []any{
			map[string]any{"role": "admin", "members": []any{"1", "2"}},
		},},
	}
	patches := []*Patch{
		{Ope: "merge", PathKey: "/spec/bindings", ByKey: "role",Value: []any{map[string]any{"role": "admin", "members": []any{"3"}},},},
	}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].
		(map[string]any)["bindings"].
		([]any)[0].
		(map[string]any)["members"].
		([]any)) != 3 {
	t.Fatalf("byKey merge failed, got %v",
			tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any), //
		)
	}
}

func TestApply_ScalarArray_Add(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"p": []any{321}},
	}
	patches := []*Patch{
			{Ope: "merge", PathKey: "/spec/p", ItemOps: "add",
			Value: []any{987},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["p"].([]any)) != 2 {
	t.Fatalf("add failed, got %v",
			tgt["spec"].(map[string]any)["p"].([]any), //
		)
	}
}

func TestApply_ScalarArray_Remove(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"p": []any{321, 987}},
	}
	patches := []*Patch{
			{Ope: "merge", PathKey: "/spec/p", ItemOps: "remove",
			Value: []any{321},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
			t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["p"].([]any)) != 1 || //
	tgt["spec"].(map[string]any)["p"].([]any)[0] != 987 ||
	false {
	t.Fatalf("remove failed, got %v",
			tgt["spec"].(map[string]any)["p"].([]any),
		)
	}
}

func TestApply_ScalarArray_Replace(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"p": []any{321, 987}},
	}
	patches := []*Patch{
			{
			Ope: "merge", PathKey: "/spec/p", ItemOps: "replace",
			Old: 321,
			Value: []any{123},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if tgt["spec"].(map[string]any)["p"].([]any)[0] != 123 || //
	tgt["spec"].(map[string]any)["p"].([]any)[1] != 987 ||
	false {
	t.Fatalf("replace failed, got %v",
			tgt["spec"].(map[string]any)["p"].([]any),
		)
	}
}

func TestApply_ScalarArray_Disable(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"p": []any{321, 987}},
	}
	patches := []*Patch{
			{Ope: "merge", PathKey: "/spec/p", ItemOps: "disable",
			Value: []any{321},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
			t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["p"].([]any)) != 1 || //
	tgt["spec"].(map[string]any)["p"].([]any)[0] != 987 ||
	false {
	t.Fatalf("disable failed, got %v",
			tgt["spec"].(map[string]any)["p"].([]any),
		)
	}
}

func TestApply_ScalarArray_Keep(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"p": []any{321, 987, 123}},
	}
	patches := []*Patch{
			{Ope: "merge", PathKey: "/spec/p", ItemOps: "keep",
			Value: []any{321, 987},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["p"].([]any)) != 2 {
	t.Fatalf("keep failed, got %v",
			tgt["spec"].(map[string]any)["p"].([]any), //
		)
	}
}

func TestApply_MixedArray_Replace(t *testing.T) {
	src := map[string]any{"spec":
	map[string]any{"tags": []any{1, "a", true}},
	}
	patches := []*Patch{
			{Ope: "replace", PathKey: "/spec/tags", MixedAr: true,
			Value: []any{"x", "y"},
			},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["tags"].([]any)) != 2 || //
	tgt["spec"].(map[string]any)["tags"].([]any)[0] != "x" ||
	false {
	t.Fatalf("mixed replace failed, got %v",
			tgt["spec"].(map[string]any)["tags"].([]any),
		)
	}
}

func TestItemOps_Expr_And(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{
					"role": "admin",
					"env":  "prod",
					"members": []any{1, 2, 3},
				},
				map[string]any{
					"role": "admin",
					"env":  "dev",
					"members": []any{4, 5},
				},
			},
		},
	}

	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role==admin && env==prod/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	members := tgt["spec"].
		(map[string]any)["bindings"].([]any)[0].
		(map[string]any)["members"].([]any)

	if len(members) != 2 {
		t.Fatalf("remove failed, got %v", members)
	}
}

// func

func TestItemOps_Expr_NotUnique(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin"},
				map[string]any{"role": "admin"},
			},
		},
	}

	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role==admin/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}

	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for non-unique expr")
	}
}

func TestItemOps_Expr_NotEqual(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{
					"role": "admin",
					"members": []any{1, 2},
				},
				map[string]any{
					"role": "viewer",
					"members": []any{3, 4},
				},
			},
		},
	}

	tgt := DeepCopy(src).(map[string]any)

	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: "/spec/bindings/role!=viewer/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}

	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}

	members := tgt["spec"].
		(map[string]any)["bindings"].([]any)[0].
		(map[string]any)["members"].([]any)

	if len(members) != 1 {
		t.Fatalf("remove with != failed, got %v", members)
	}
}

func TestApply_Delete_FromStructArray(t *testing.T) {
	src := map[string]any{"spec":
	//
	map[string]any{"bindings":
	//
	[]any{
	//
	map[string]any{"role": "a"},
	//
	map[string]any{"role": "b"},
	//
	},
	//
	},
	//
	}
	patches := []*Patch{
			{Ope: "delete", PathKey: "/spec/bindings", ByKey: "role", Value: "a"}, ///**
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["bindings"].([]any)) != 1 { //
		t.Fatalf("delete from struct array failed")
	}
}

func TestApply_Delete_Field(t *testing.T) {
	src := map[string]any{"spec": map[string]any{"add": "r"}} //
	patches := []*Patch{
			{Ope: "delete", PathKey: "/spec/add"},
		}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if _, ok := tgt["spec"].(map[string]any)["add"]; ok {
		t.Fatalf("delete field failed")
	}
}

func init() {
	///**
}

// struct

// interface