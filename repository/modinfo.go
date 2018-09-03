package repository

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

type ModInfo struct {
	Source  string `json:"source"`
	Version string `json:"version"`
}

func (mi ModInfo) String() string {
	return fmt.Sprintf("(%s @ %s)", mi.Source, mi.Version)
}

var (
	// examples
	//  mod file style
	//   github.com/foo/bar v2.0.0
	//   github.com/tdewolff/parse v2.3.3+incompatible // indirect
	//   golang.org/x/tools v0.0.0-20180111040409-fbec762f837d
	//   gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405
	//  proxy request style
	//   /github.com/cpuguy83/go-md2man/@v/v1.0.6.info
	modLineRe = regexp.MustCompile(`([\S.-]+)([\s]+)(v[\d]+.[\d]+.[\d]+(-[\d]+-[0-9a-f]+)?)`)
)

// Parse will parse s as a module in string form.
func Parse(s string) (ModInfo, error) {
	s = strings.Trim(s, "/")
	s = strings.Replace(s, "/@v/", " ", -1)
	groups := modLineRe.FindStringSubmatch(s)

	fmt.Println("xx groups:", len(groups), groups)
	for i := 0; i < len(groups); i++ {
		fmt.Printf("groups[%d]: %s\n", i, groups[i])
	}

	var mod ModInfo
	if len(groups) != 5 {
		return mod, errors.Errorf("malformed module line: %q", s)
	}
	return ModInfo{
		Source:  strings.TrimRight(groups[1], "/"),
		Version: groups[3],
	}, nil
}
