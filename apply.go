package smpatch

import "fmt"

func Apply(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
) error {

	tgt = DeepCopy(src).(map[string]any)

	for _, p := range patches {
		switch {
		case p.Ope == "delete":
			if err := applyDel(p, tgt); err != nil {
				return err
			}

		case p.Ope == "replace" && p.MixedAr:
			if err := mixed(p, tgt); err != nil {
				return err
			}

		case p.Ope == "replace":
			if err := applyRpl(p, tgt); err != nil {
				return err
			}

		case p.Ope == "merge" && p.ByKey != "":
			if err := mapAr(p, tgt); err != nil {
				return err
			}

		case p.Ope == "merge" && p.ItemOps != "":
			if err := itemOps(p, tgt); err != nil {
				return err
			}

		case p.Ope == "merge":
			if err := dpMeg(p, tgt); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown ope: %s", p.Ope)
		}
	}
	return nil
}