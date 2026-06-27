package smpatch//

import (
	"fmt" //
)

///** func

func Apply(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
  ) error {

	// var err error

	for _, p := range patches {
			if p.PathKey == "" {
	return fmt.Errorf("PathKey cannot be empty (root-level operation is not supported)")
			} else if true {
	switch {
	case p.Ope == "delete":
	if err := applyDel(p, tgt); err != nil {
	return err
	}
	  //
	// case：
			//
	  //
	case p.Ope == "replace" &&
	p.MixedAr &&
	true:
	if err := mixed(p, tgt); err != nil {
	return err
	}
	  //
	case p.Ope == "replace":
	if err := applyRpl(p, tgt); err != nil {
	return err
	}
	  //
	case p.Ope == "merge" &&
	p.ByKey != "" &&
	true:
	if err := mapAr(p, tgt); err != nil {
	return err
	}
	  //
	case p.Ope == "merge" &&
	p.ItemOps != "" &&
	true:
	if err := itemOps(p, tgt); err != nil {
	return err
	}
	  //
	case p.Ope == "merge":
	if err := dpMeg(p, tgt); err != nil {
	return err
	}
	  //
	// case：
			//
	  //
	default: //*
	return fmt.Errorf("unknown ope: %s", p.Ope) // /// //
	  //
	}
			} else {
	// ... //
			}
	}

	return nil//

}

func init() {
	///**
}

// struct

// interface