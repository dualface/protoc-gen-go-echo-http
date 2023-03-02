// Package annotation provides a parser for annotations defined with "@" in comments
package annotation

import (
	"regexp"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

type (
	// Annotation Used to save annotations defined with "@" in comments
	// Syntax:
	// @Name(Key1=V, Key2=V, ...)
	Annotation struct {
		Name  string
		Value map[string]string
	}
)

// ParseAnnotations Parse annotations from comments
func ParseAnnotations(cs protogen.CommentSet) []Annotation {
	as := []Annotation{}
	for _, s := range cs.LeadingDetached {
		a := parseAnnotation(s.String())
		if a != nil {
			as = append(as, *a)
		}
	}
	{
		a := parseAnnotation(cs.Leading.String())
		if a != nil {
			as = append(as, *a)
		}
	}
	{
		a := parseAnnotation(cs.Trailing.String())
		if a != nil {
			as = append(as, *a)
		}
	}
	return as
}

var (
	matchAnnotation = regexp.MustCompile(`@[ \t]*([^ \t\(\)]+)[ \t]*(\(.+\))?.*`)
	matchValue      = regexp.MustCompile(`([^ \t]+)[ \t]*=[ \t]*([^ \t]+)`)
)

func parseAnnotation(line string) *Annotation {
	line = strings.Trim(line, "/ \t\n")
	m := matchAnnotation.FindStringSubmatch(line)
	if len(m) < 2 {
		return nil
	}

	a := &Annotation{
		Name: strings.ToLower(m[1]),
	}

	if len(m) >= 3 {
		args := strings.Split(strings.Trim(m[2], " \t()"), ",")
		a.Value = make(map[string]string)
		for _, arg := range args {
			p := matchValue.FindAllStringSubmatch(arg, -1)
			for _, pair := range p {
				a.Value[strings.Trim(pair[1], `"'`)] = strings.Trim(pair[2], `"'`)
			}
		}
	}

	return a
}
