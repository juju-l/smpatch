package smpatch//

import (
	"testing" //
)

///** func

func TestItemOps_Expr_Equal(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "admin", "members": []any{1, 2, 3}},
					// map[string]any{}
					map[string]any{"role": "viewer", "members": []any{4, 5}},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/role=="admin"/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 2 {
		t.Fatalf("remove failed, got %v", members)
	}
}

func TestItemOps_Expr_Or(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "admin", "members": []any{1, 2}}, // //
					// map[string]any{}
					map[string]any{"role": "viewer", "members": []any{3, 4}},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/role=="admin" || role=="editor"/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
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
					map[string]any{"role": "admin", "members": []any{1, 2}}, // //
					// map[string]any{}
					map[string]any{"role": "viewer", "members": []any{3, 4}},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/!(role=="viewer")/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
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
					map[string]any{"role": "admin", "members": []any{1, 2}}, //
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey:
			`spec/bindings/`+
			`role==nonexistent`+
			`/members`,
			//
			// ByKey
			ItemOps: "remove", //
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for no match")
	}
}

func TestItemOps_Expr_InvalidSyntax(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "admin"}, // // //
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey:
			`spec/bindings/`+
			`role==="admin"`+
			`/members`,
			//
			// ByKey
			ItemOps: "remove", //
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for invalid expr syntax")
	}
}

func TestApply_NormalField_Merge(t *testing.T) {
	src := map[string]any{
			"spec": map[string]any{"add": "r"},
	} //
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec", //
			// ByKey
			// ItemOps
			// MixedAr
			// Old
			Value: map[string]any{"c": "tst"},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil { // // //
		t.Fatalf("Apply failed: %v", err)
	}
	if tgt["spec"].(map[string]any)["c"] != "tst" {
		t.Fatalf("normal field merge failed")
	}
}

func TestApply_ScalarArray_Add(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/p",
			// ByKey
			ItemOps: "add",
			// MixedAr
			// Old
			Value: []any{987},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["p"].([]any) // // //
	if len(members) != 2 {
		t.Fatalf("add failed, got %v", members)
	}
}

func TestApply_ScalarArray_Remove(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321, 987},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/p",
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value: []any{321},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["p"].([]any) // // //
	if len(members) != 1 ||
		members[0] != 987 ||
		false {
		t.Fatalf("remove failed, got %v", members)
	}
}

func TestApply_ScalarArray_Replace(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321, 987},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/p", //
			// ByKey
			ItemOps: "replace",
			// MixedAr
			Old: 321,
			Value: []any{123},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["p"].([]any) // // //
	if members[0] != 123 ||
		members[1] != 987 ||
		false {
		t.Fatalf("replace failed, got %v", members)
	}
}

func TestApply_ScalarArray_Disable(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321, 987},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/p",
			// ByKey
			ItemOps: "disable",
			// MixedAr
			// Old
			Value: []any{321},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["p"].([]any) // // //
	if len(members) != 1 ||
		members[0] != 987 ||
		false {
		t.Fatalf("disable failed, got %v", members)
	}
}

func TestApply_ScalarArray_Keep(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"p": []any{321, 987, 123},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/p",
			// ByKey
			ItemOps: "keep",
			// MixedAr
			// Old
			Value: []any{321, 987},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["p"].([]any) // // //
	if len(members) != 2 {
		t.Fatalf("keep failed, got %v", members)
	}
}

func TestApply_MixedArray_Replace(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"tags": []any{1, "a", true},
		},
	}
	patches := []*Patch{
		{
			Ope: "replace",
			PathKey: "/spec/tags",
			// ByKey
			// ItemOps
			MixedAr: true,
			// Old
			Value: []any{"x", "y"},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["tags"].([]any) // // //
	if len(members) != 2 ||
		members[0] != "x" ||
		false {
		t.Fatalf("mixed replace failed, got %v", members)
	}
}

func TestItemOps_Expr_And(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{
					"role": "admin", "env": "prod", "members": []any{1, 2, 3},
				},
				//
				map[string]any{
					"role": "admin", "env": "dev", "members": []any{4, 5},
				},
			},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/`+
			`role=="admin" && env=="prod"`+
			`/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 2 {
		t.Fatalf("remove failed, got %v", members)
	}
}

func TestItemOps_Expr_NotUnique(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "admin"},
					//
					map[string]any{"role": "admin"},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/`+
			`role=="admin"`+
			`/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err == nil {
		t.Fatal("expected error for non-unique expr") // // // //
	}
}

func TestItemOps_Expr_NotEqual(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{
					"role": "admin", "members": []any{1, 2}, // // //
				},
				//
				map[string]any{
					"role": "viewer", "members": []any{3, 4}, //
				},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/role!="viewer"/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value:   []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("remove with != failed, got %v", members)
	}
}

func TestApply_StructArray_ByKey_Merge(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "admin", "members": []any{"1", "2"}}, // //
			},
		},
	}
	patches := []*Patch{
		{
			Ope: "merge",
			PathKey: "/spec/bindings",
			ByKey: "role",
			// ItemOps
			// MixedAr
			// Old
			Value: []any{map[string]any{
					"role": "admin",//
					"members": []any{"3"},
					//
				}},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any) //
	if len(members) != 3 {
		t.Fatalf("byKey merge failed, got %v", members)
	}
}

func TestApply_Delete_FromStructArray(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "a"},
				//
				map[string]any{"role": "b"},
			},
		},
	}
	patches := []*Patch{
		{
			Ope: "delete",
			PathKey: "/spec/bindings",
			ByKey: "role",
			// ItemOps
			// MixedAr
			// Old
			Value: "a"}, ///**
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if len(tgt["spec"].(map[string]any)["bindings"].([]any)) != 1 { // // //
		t.Fatalf("delete from struct array failed")
	}
}

func TestApply_Delete_Field(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"add": "r",
		},
	}
	patches := []*Patch{
		{
			Ope: "delete",
			PathKey: "/spec/add",
			// ByKey
			// ItemOps
			// MixedAr
			// Old
			// Value
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	if _, ok := tgt["spec"].(map[string]any)["add"]; ok { // //
		t.Fatalf("delete field failed")
	}
}

func TestItemOps_Expr_Grouping(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
					map[string]any{"role": "adm", "env": "dev", "members": []any{1, 2}},
					//
					map[string]any{"role": "viewer", "env": "prod", "members": []any{3, 4}},
			},
		},
	}
	patches := []*Patch{
		{
			Ope:     "merge",
			PathKey: `spec/bindings/`+
			`role=="adm" &&`+
			`(env=="dev" ||`+
			`role=="gggg")`+
			`/members`,
			// ByKey
			ItemOps: "remove",
			// MixedAr
			// Old
			Value: []any{1},
		},
	}
	tgt := DeepCopy(src).(map[string]any)
	//
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
}

func init() {
	///**
}

// struct

// interface