package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestRead(t *testing.T) {
	files := []string{"main_test.go", "README.md", "main.go", "LICENSE", "go.mod"}

	current := Read(".")

	s := make(map[string]bool)

	for _, i := range current {
		if strings.HasPrefix(i, ".") {
			continue
		}

		for _, f := range files {
			if f == i {
				s[i] = true
			}
		}
	}

	if len(s) != len(files) {
		t.Errorf("invalid number of files: %v\nfiles:\n", len(s))
		for s, _ := range s {
			fmt.Println(s)
		}
	}
}
