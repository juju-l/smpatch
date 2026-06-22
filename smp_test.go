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

func TestItemOps_Expr_NotEqual(t *testing.T) {
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
			PathKey: "/spec/bindings/role!=viewer/members",
			ItemOps: "remove",
			Value:   []any{1},
		},
	}
	if err := Apply(src, patches, tgt); err != nil {
		t.Fatalf("Apply failed: %v", err)
	}
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 1 {
		t.Fatalf("remove with != failed, got %v", members)
	}
}

func TestItemOps_Expr_And(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "env": "prod", "members": []any{1, 2, 3}},
				map[string]any{"role": "admin", "env": "dev", "members": []any{4, 5}},
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
	members := tgt["spec"].(map[string]any)["bindings"].([]any)[0].(map[string]any)["members"].([]any)
	if len(members) != 2 {
		t.Fatalf("remove with && failed, got %v", members)
	}
}

func TestItemOps_Expr_Or(t *testing.T) {
	src := map[string]any{
		"spec": map[string]any{
			"bindings": []any{
				map[string]any{"role": "admin", "members": []any{1, 2}},
				map[string]any{"role": "editor", "members": []any{3, 4}},
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
		t.Fatal("expected error for non-unique expr result")
	}
}

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