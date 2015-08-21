package messagediff

import (
	"errors"
	"reflect"
)

var errPtrRequired = errors.New("can only patch referenced types (pointers)")

// Patch applies a diff to a struct.
func Patch(a interface{}, diff *Diff) (interface{}, error) {
	if reflect.TypeOf(a).Kind() != reflect.Ptr {
		return nil, errPtrRequired
	}
	var err error
	for path, insert := range diff.Added {
		if a, err = addMod(a, *path, insert); err != nil {
			return nil, err
		}
	}
	for path, mod := range diff.Modified {
		if a, err = addMod(a, *path, mod); err != nil {
			return nil, err
		}
	}
	for path, mod := range diff.Removed {
		if a, err = rem(a, *path, mod); err != nil {
			return nil, err
		}
	}
	return a, nil
}

func addMod(a interface{}, path Path, change interface{}) (interface{}, error) {
	cv := reflect.ValueOf(change)
	v := reflect.ValueOf(a).Elem()
	for i, n := range path {
		isLast := i == len(path)-1
		switch np := n.(type) {
		case SliceIndex:
			si := int(np)
			if isLast {
				if v.Len() < si+1 {
					zero := reflect.Zero(cv.Type())
					for j := 0; j <= si+1-v.Len(); j++ {
						v.Set(reflect.Append(v, zero))
					}
				}
				v.Index(si).Set(cv)
			} else {
				v = v.Index(si)
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
			if !v.CanSet() {
				v = unsafeReflectValue(v)
			}
			if isLast {
				v.Set(cv)
			}
		default:
			panic("unknown StructField")
		}
	}
	return a, nil
}

func rem(a interface{}, path Path, change interface{}) (interface{}, error) {
	v := reflect.ValueOf(a).Elem()
	for i, n := range path {
		isLast := i == len(path)-1
		switch np := n.(type) {
		case SliceIndex:
			si := int(np)
			if isLast {
				for j := min(v.Len()-1, si); j >= 0; j-- {
					if reflect.DeepEqual(v.Index(j).Interface(), change) {
						v.Set(reflect.AppendSlice(v.Slice(0, j), v.Slice(j+1, v.Len())))
						break
					}
				}
			} else {
				v = v.Index(si)
			}
		case MapKey:
			key := reflect.ValueOf(np.Key)
			if isLast {
				v.SetMapIndex(key, reflect.Value{})
			} else {
				v = v.MapIndex(key)
			}
		case StructField:
			v = v.FieldByName(string(np))
			if isLast {
				v.Set(reflect.Zero(v.Type()))
			}
		}
	}
	return a, nil
}
