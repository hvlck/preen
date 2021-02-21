package preen

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

func TestFind(t *testing.T) {
	u := "https://github.com\nhttps://bit.ly"
	urls := strings.ReplaceAll(u, "\r", "")

	matches := Find(urls)
	fmt.Println(matches)

	if len(matches) != 2 {
		t.Errorf("found invalid number of urls: %v when length should be 2", len(matches))
	}
}
