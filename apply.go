package smpatch

import "fmt"

func Apply(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
) error {
	for _, p := range patches {
		switch {
		case p.Ope == "delete":
			if err := applyDel(p, src, tgt); err != nil {
				return err
			}

		case p.Ope == "replace" && p.MixedAr:
			if err := mixed(p, src, tgt); err != nil {
				return err
			}

		case p.Ope == "replace":
			if err := applyRpl(p, src, tgt); err != nil {
				return err
			}

		case p.Ope == "merge" && p.ByKey != "":
			if err := mapAr(p, src, tgt); err != nil {
				return err
			}

		case p.Ope == "merge" && p.ItemOps != "":
			if err := itemOps(p, src, tgt); err != nil {
				return err
			}

		case p.Ope == "merge":
			if err := dpMeg(p, src, tgt); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown ope: %s", p.Ope)
		}
	}
	return nil
}