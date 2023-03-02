// go test -v annotation_test.go annotation.go
package annotation

import (
	"testing"
)

func TestParseAnnotationFailue(t *testing.T) {
	ss := []string{
		"//Annotation",
		"//@",
	}

	for _, s := range ss {
		a := parseAnnotation(s)
		if a != nil {
			t.Fatalf("parse `%s` failed", s)
		}
	}
}

func TestParseAnnotationWithoutArguments(t *testing.T) {
	ss := []string{
		"//@Annotation",
		"// @Annotation",
		"//\t@Annotation",
		"// \t@Annotation",
		"//@ Annotation",
		"//@\tAnnotation",
		"//@ \tAnnotation",
		"// @ \tAnnotation",
		"//\t@ \tAnnotation",
		"// \t@ \tAnnotation",
		"// @Annotation()",
		"// @Annotation ()",
		"// @Annotation \t()",
		"// @Annotation ( )",
		"// @Annotation \t( \t)",
	}

	for _, s := range ss {
		a := parseAnnotation(s)
		if a == nil {
			t.Fatalf("parse `%s` failed", s)
		}

		if a.Name != "annotation" {
			t.Fatalf("expected Name is `annotation`, actual is `%s`", a.Name)
		}
		if len(a.Value) != 0 {
			t.Fatalf("expected Value is empty, actual has %d elements", len(a.Value))
		}

		t.Logf("parse ok, %#v", a)
	}
}

func TestParseAnnotationWithArguments(t *testing.T) {
	ss := []string{
		"//@Annotation(k1=v1)",
		"//@Annotation(k1=v1,k2=v2)",
		"//@Annotation(k1=v1, k2=v2,k3=v3)",
		"//@Annotation(k1=v1, \tk2=v2,   k3=v3,\t\tk4=v4.ext)",
	}

	for i, s := range ss {
		a := parseAnnotation(s)
		if a == nil {
			t.Fatalf("parse `%s` failed", s)
		}

		if a.Name != "annotation" {
			t.Fatalf("expected Name is `annotation`, actual is `%s`", a.Name)
		}
		c := i + 1
		if len(a.Value) != c {
			t.Fatalf("expected Value has %d elements, actual is empty, line is `%s`", c, s)
		}

		t.Logf("parse ok, %#v", a)
	}
}
