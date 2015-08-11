package messagediff

import (
	"fmt"
	"reflect"
)

// DeepDiff does a deep comparison and returns the results.
func DeepDiff(a, b interface{}) (string, bool) {
	return diff(a, b, "")
}

func diff(a, b interface{}, path string) (string, bool) {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Type() != bVal.Type() {
		return fmt.Sprintf("modified: %s = %#v\n", path, b), false
	}
	kind := aVal.Type().Kind()
	switch kind {
	case reflect.Array, reflect.Slice:
		var cDiff string
		aLen := aVal.Len()
		bLen := bVal.Len()
		for i := 0; i < min(aLen, bLen); i++ {
			localPath := fmt.Sprintf("%s[%d]", path, i)
			d, equal := diff(aVal.Index(i).Interface(), bVal.Index(i).Interface(), localPath)
			if equal {
				continue
			}
			cDiff += d
		}
		if aLen > bLen {
			for i := bLen; i < aLen; i++ {
				localPath := fmt.Sprintf("%s[%d]", path, i)
				cDiff += fmt.Sprintf("removed: %s = %#v\n", localPath, aVal.Index(i).Interface())
			}
		} else if aLen < bLen {
			for i := aLen; i < bLen; i++ {
				localPath := fmt.Sprintf("%s[%d]", path, i)
				cDiff += fmt.Sprintf("added: %s = %#v\n", localPath, bVal.Index(i).Interface())
			}
		}
		return cDiff, len(cDiff) == 0
	case reflect.Map:
		return "TODO(d4l3k): Maps", false
	case reflect.Struct:
		return "TODO(d4l3k): Structs", false
	default:
		if reflect.DeepEqual(a, b) {
			return "", true
		}
		return fmt.Sprintf("modified: %s = %#v\n", path, b), false
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
