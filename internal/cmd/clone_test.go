package cmd

import (
	"fmt"
	"testing"
)

func TestGetRepoPage(t *testing.T) {
	s, err := getRepoPage("github.com/neovim", 1)
	if err != nil {
		t.Error(err)
	}
	if s.maxPage != 2 {
		t.Error("expected max page to be 2")
	}
	fmt.Println(s)
}
