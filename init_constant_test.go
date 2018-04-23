package main

import "testing"

func TestConstant(t *testing.T) {
	if len(strategies) != 28 {
		t.Error("-- strategies expected to be 28")
	}
}
