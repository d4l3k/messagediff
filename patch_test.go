package messagediff

import "testing"

func TestPatch(t *testing.T) {
	testData := []struct {
		a, b interface{}
	}{{
		&[]int{1},
		&[]int{1, 2, 3},
	}}
	for i, td := range testData {
		diff, equal := DeepDiff(td.a, td.b)
		if equal {
			t.Errorf("%d. DeepDiff(%#v, %#v) equal = %#v; not %#v", i, td.a, td.b, equal, false)
		}
		out := Patch(td.a, diff)
		if d, eq := PrettyDiff(td.b, out); !eq {
			t.Errorf("%d. Patch(%#v, %#v) equal = %#v; diff %s", i, td.a, diff, out, d)
		}
	}
}
