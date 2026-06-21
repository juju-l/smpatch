package smpatch//

// import (
// 	
// )

///** func

func DiffCmp(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
  ) error {
	return DyRun(src, patches, tgt) ///**
}

func init() {
	///**
}

// struct

// interface