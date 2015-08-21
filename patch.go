package messagediff

import (
	"log"
	"reflect"
)

// Patch applies a diff to a struct.
func Patch(a interface{}, diff *Diff) interface{} {
	for path, insert := range diff.Added {
		a = addMod(a, *path, insert)
	}
	for path, mod := range diff.Modified {
		a = addMod(a, *path, mod)
	}
	for path, mod := range diff.Removed {
		a = rem(a, *path, mod)
	}
	return a
}

func addMod(a interface{}, path Path, change interface{}) interface{} {
	cv := reflect.ValueOf(change)
	v := reflect.ValueOf(a)
	for i, n := range path {
		isLast := i == len(path)-1
		switch np := n.(type) {
		case SliceIndex:
			si := int(np)
			if isLast {
				log.Printf("SETLEN %d %s %#v", si, v.Kind().String(), v.Interface())
				if v.Len() >= si {
					v.SetLen(si + 1)
				}
			}
			v = v.Index(si)
			if isLast {
				v.Set(cv)
			}
		case MapKey:
			key := reflect.ValueOf(np.Key)
			if isLast {
				v.SetMapIndex(key, cv)
			} else {
				v = v.MapIndex(key)
			}
		case StructField:
			v = v.FieldByName(string(np))
			if isLast {
				v.Set(cv)
			}
		}
	}
	return a
}
func rem(a interface{}, path Path, change interface{}) interface{} {
	return a
}
