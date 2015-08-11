# messagediff
A library for doing diffs of arbitrary Golang structs.

If the unsafe package is available messagediff will diff unexported fields in
addition to exported fields. This is primarily used for testing purposes as it
allows for providing informative error messages.


## Example Usage
In a normal file:
```go
package main

import "github.com/d4l3k/messagediff"

type someStruct struct {
  A, b int
  C []int
}

func main() {
			a := someStruct{1, 2, []int{1}}
			b := someStruct{1, 3, []int{1, 2}}
      diff, equal := messagediff.DeepDiff(a, b)
      /*
        diff =
          `added: .C[1] = 2
          modified: .b = 3`

        equal = false
      */
}

```
In a test:
```go
import "github.com/d4l3k/messagediff"

...

type someStruct struct {
  A, b int
  C []int
}

func TestSomething(t *testing.T) {
  want := someStruct{1, 2, []int{1}}
  got := someStruct{1, 3, []int{1, 2}}
  if diff, equal := messagediff.DeepDiff(want, got); !equal {
    t.Errorf("Something() = %#v\n%s", got, diff)
  }
}
```

## Copyright
Copyright (c) 2015 Tristan Rice <rice@fn.lc>

messagediff is licensed under the MIT license. See the LICENSE file for more information.

bypass.go and bypasssafe.go are borrowed from
[go-spew](https://github.com/davecgh/go-spew) and have a seperate copyright
notice.
