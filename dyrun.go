package smpatch//

// import (
// 	
// )

///** func

func DyRun(
	src map[string]any,
	patches []*Patch,
	tgt map[string]any,
  ) error {
	return Apply(src, patches, tgt) ///**
}

func init() {
	///**
}

// struct

// interface