package messagediff

import "testing"

type testStruct struct {
	A, b int
	C    []int
}

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
			map[string]int{"b": 4, "c": 3},
			"added: [\"c\"] = 3\nmodified: [\"b\"] = 4\nremoved: [\"a\"] = 1\n",
			false,
		},
		{
			testStruct{1, 2, []int{1}},
			testStruct{1, 3, []int{1, 2}},
			"added: .C[1] = 2\nmodified: .b = 3\n",
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
