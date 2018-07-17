package xstrings

import (
	"fmt"
	"strings"
	"testing"
)

type testcase struct {
	result string
	sep    string
	args   []interface{}
}

func (tc *testcase) Check(t *testing.T) {
	s := Join(tc.args, tc.sep)
	if s != tc.result {
		t.Errorf("expected %v got %v", tc.result, s)
	}
}

func TestJoinStrings(t *testing.T) {
	tc := testcase{
		result: "abc",
		sep:    "",
		args:   []interface{}{"a", "b", "c"},
	}
	tc.Check(t)

	tc.sep = "-"
	tc.result = "a-b-c"
	tc.Check(t)
}

func TestJoinMixedPrimitive(t *testing.T) {
	tc := testcase{
		result: "a2c4.33",
		sep:    "",
		args:   []interface{}{"a", 2, "c", 4.33},
	}
	tc.Check(t)
}

type foo struct {
	bar int
	baz []string
}

func (f foo) String() string {
	return fmt.Sprintf("%v->%v", f.bar, strings.Join(f.baz, "_"))
}

func TestJoinMixed(t *testing.T) {
	tc := testcase{
		result: "a/5/2->abc_xyz",
		sep:    "/",
		args: []interface{}{
			"a", 5, foo{
				bar: 2,
				baz: []string{"abc", "xyz"},
			},
		},
	}
	tc.Check(t)
}
