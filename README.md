# enumgen

A utility to generate reverse-lookup tables for enumerated types in Go.

## Usage

```bash
Usage of enumgen:
        enumgen [flags] -type T [directory]
  -output string
        output file name; default srcdir/<type>_enum.go
  -type string
        comma-separated list of type names; must be set
```

## Example

Given the following declaration and `enumgen` invocation:

```go
package foo

//go:generate enumgen -type=MyEnum
type MyEnum int

const (
    A MyEnum = iota
    B
    C
)
```

The generated file will be named `myenum_enum.go`, and will have the following content:

```go
package foo

var _MyEnumValues = map[int]MyEnum{
        0: A,
        1: B,
        2: C,
}

func GetMyEnum(x int) (MyEnum, bool) {
        v, ok := _MyEnumValues[x]
        return v, ok
}
```
