package messagediff

import "testing"

func TestDeepDiff(t *testing.T) {
	testData := []struct {
		a, b  interface{}
		diff  string
		equal bool
	}{
		{
			true,
			false,
			"modified:  = false\n",
			false,
		},
		{
			true,
			0,
			"modified:  = 0\n",
			false,
		},
		{
			[]int{0, 1, 2},
			[]int{0, 1, 2, 3},
			"added: [3] = 3\n",
			false,
		},
		{
			[]int{0, 1, 2, 3},
			[]int{0, 1, 2},
			"removed: [3] = 3\n",
			false,
		},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"b": 2, "c": 3},
			"removed: [\"a\"] = 1\nadded: [\"c\"] = 3\n",
			false,
		},
	}
	for i, td := range testData {
		diff, equal := DeepDiff(td.a, td.b)
		if diff != td.diff {
			t.Errorf("%d. DeepDiff(%#v, %#v) diff = %#v; not %#v", i, td.a, td.b, diff, td.diff)
		}
		if equal != td.equal {
			t.Errorf("%d. DeepDiff(%#v, %#v) equal = %#v; not %#v", i, td.a, td.b, equal, td.equal)
		}
	}
}
