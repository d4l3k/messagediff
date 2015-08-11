package messagediff

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// PrettyDiff does a deep comparison and returns the results.
func PrettyDiff(a, b interface{}) (string, bool) {
	d, equal := diff(a, b, "")
	sort.Strings(d)
	return strings.Join(d, ""), equal
}

func diff(a, b interface{}, path string) ([]string, bool) {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	if aVal.Type() != bVal.Type() {
		return []string{fmt.Sprintf("modified: %s = %#v\n", path, b)}, false
	}
	kind := aVal.Type().Kind()
	switch kind {
	case reflect.Array, reflect.Slice:
		var cDiff []string
		aLen := aVal.Len()
		bLen := bVal.Len()
		for i := 0; i < min(aLen, bLen); i++ {
			localPath := fmt.Sprintf("%s[%d]", path, i)
			d, equal := diff(aVal.Index(i).Interface(), bVal.Index(i).Interface(), localPath)
			if equal {
				continue
			}
			cDiff = append(cDiff, d...)
		}
		if aLen > bLen {
			for i := bLen; i < aLen; i++ {
				localPath := fmt.Sprintf("%s[%d]", path, i)
				cDiff = append(cDiff, fmt.Sprintf("removed: %s = %#v\n", localPath, aVal.Index(i).Interface()))
			}
		} else if aLen < bLen {
			for i := aLen; i < bLen; i++ {
				localPath := fmt.Sprintf("%s[%d]", path, i)
				cDiff = append(cDiff, fmt.Sprintf("added: %s = %#v\n", localPath, bVal.Index(i).Interface()))
			}
		}
		return cDiff, len(cDiff) == 0
	case reflect.Map:
		var cDiff []string
		for _, key := range aVal.MapKeys() {
			aI := aVal.MapIndex(key)
			bI := bVal.MapIndex(key)
			localPath := fmt.Sprintf("%s[%#v]", path, key.Interface())
			if !bI.IsValid() {
				cDiff = append(cDiff, fmt.Sprintf("removed: %s = %#v\n", localPath, aI.Interface()))
			} else if d, equal := diff(aI.Interface(), bI.Interface(), localPath); !equal {
				cDiff = append(cDiff, d...)
			}
		}
		for _, key := range bVal.MapKeys() {
			aI := aVal.MapIndex(key)
			bI := bVal.MapIndex(key)
			localPath := fmt.Sprintf("%s[%#v]", path, key.Interface())
			if !aI.IsValid() {
				cDiff = append(cDiff, fmt.Sprintf("added: %s = %#v\n", localPath, bI.Interface()))
			}
		}
		return cDiff, len(cDiff) == 0
	case reflect.Struct:
		var cDiff []string
		typ := aVal.Type()
		for i := 0; i < typ.NumField(); i++ {
			index := []int{i}
			field := typ.FieldByIndex(index)
			localPath := fmt.Sprintf("%s.%s", path, field.Name)
			aI := unsafeReflectValue(aVal.FieldByIndex(index)).Interface()
			bI := unsafeReflectValue(bVal.FieldByIndex(index)).Interface()
			if d, equal := diff(aI, bI, localPath); !equal {
				cDiff = append(cDiff, d...)
			}
		}
		return cDiff, len(cDiff) == 0
	default:
		if reflect.DeepEqual(a, b) {
			return nil, true
		}
		return []string{fmt.Sprintf("modified: %s = %#v\n", path, b)}, false
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
