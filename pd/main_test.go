package main

import "testing"

func TestInvokeCLIVersion(t *testing.T) {
	args := []string{"-v"}
	if invokeCLI(args) != 0 {
		t.Errorf("`pd -v` return code is non zero")
	}
}
