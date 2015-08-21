package messagediff

import (
	"errors"
	"reflect"
)

var (
	errPtrRequired     = errors.New("can only patch referenced types (pointers)")
	errUnknownPathNode = errors.New("unknown PathNode in Path")
)

// Patch applies a diff to a struct.
func Patch(a interface{}, diff *Diff) error {
	if reflect.TypeOf(a).Kind() != reflect.Ptr {
		return errPtrRequired
	}
	v := reflect.ValueOf(a).Elem()
	for path, insert := range diff.Added {
		if err := addMod(v, *path, insert); err != nil {
			return err
		}
	}
	for path, mod := range diff.Modified {
		if err := addMod(v, *path, mod); err != nil {
			return err
		}
	}
	for path, mod := range diff.Removed {
		if err := rem(v, *path, mod); err != nil {
			return err
		}
	}
	return nil
}

func addMod(v reflect.Value, path Path, change interface{}) error {
	cv := reflect.ValueOf(change)
	for i, n := range path {
		isLast := i == len(path)-1
		switch np := n.(type) {
		case SliceIndex:
			si := int(np)
			if isLast {
				// Grow length if necessary
				nlen := si + 1
				if v.Len() < nlen {
					zero := reflect.Zero(cv.Type())
					for j := 0; j <= nlen-v.Len(); j++ {
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
			return errUnknownPathNode
		}
	}
	return nil
}

func rem(v reflect.Value, path Path, change interface{}) error {
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
		default:
			return errUnknownPathNode
		}
	}
	return nil
}
