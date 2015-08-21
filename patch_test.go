package messagediff

import "testing"

func TestPatch(t *testing.T) {
	testData := []struct {
		a, b interface{}
	}{
		{
			&[]int{1},
			&[]int{0, 2, 3},
		},
		{
			&map[string]int{
				"duck":  5,
				"mouse": 1,
			},
			&map[string]int{
				"duck": 6,
				"blue": 9,
			},
		},
		{
			&testStruct{
				5, 6, nil,
			},
			&testStruct{
				6, 7, []int{1, 2},
			},
		},
		{
			&[]int{0, 2, 3},
			&[]int{1},
		},
	}
	for i, td := range testData {
		diff, equal := DeepDiff(td.a, td.b)
		if equal {
			t.Errorf("%d. DeepDiff(%#v, %#v) equal = %#v; not %#v", i, td.a, td.b, equal, false)
		}
		if err := Patch(td.a, diff); err != nil {
			t.Errorf("%d. Patch(%#v, diff) errored %s", i, td.a, err)
		}
		if d, eq := PrettyDiff(td.b, td.a); !eq {
			t.Errorf("%d. Patch(%#v, diff) = %#v; diff %s", i, td.a, td.a, d)
		}
	}
}
