package config

import (
	"fmt"
	"reflect"
)

// getType converts string type to reflect.Type
func getType(typeName string) reflect.Type {
	switch typeName {
	case "string":
		return reflect.TypeOf("")
	case "int":
		return reflect.TypeOf(int(0))
	case "int8":
		return reflect.TypeOf(int8(0))
	case "int16":
		return reflect.TypeOf(int16(0))
	case "int32":
		return reflect.TypeOf(int32(0))
	case "int64":
		return reflect.TypeOf(int64(0))
	case "uint":
		return reflect.TypeOf(uint(0))
	case "uint8":
		return reflect.TypeOf(uint8(0))
	case "uint16":
		return reflect.TypeOf(uint16(0))
	case "uint32":
		return reflect.TypeOf(uint32(0))
	case "uint64":
		return reflect.TypeOf(uint64(0))
	case "float32":
		return reflect.TypeOf(float32(0))
	case "float64":
		return reflect.TypeOf(float64(0))
	case "bool":
		return reflect.TypeOf(false)
	default:
		if len(typeName) > 1 {
			if typeName[0] == '[' && typeName[1] == ']' {
				return reflect.SliceOf(getType(typeName[2:]))
			}
			if typeName[0] == '*' {
				return reflect.PtrTo(getType(typeName[1:]))
			}
		}
		panic(fmt.Sprintf("unsupported field type: %s", typeName))
	}
}
