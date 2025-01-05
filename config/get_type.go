package config

import (
	"fmt"
	"reflect"
	"time"
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
	case "time.Time":
		return reflect.TypeOf(time.Time{})
	case "map[string]string":
		return reflect.TypeOf(map[string]string{})
	default:
		// Check for slices
		if len(typeName) > 2 && typeName[0] == '[' && typeName[1] == ']' {
			// For []struct, this should be processed separately
			if typeName[2:6] == "struct" {
				// Remove '[]struct' from the typeName
				structType := typeName[7 : len(typeName)-1] // Removing '[]struct{...}'
				return reflect.SliceOf(getType("struct{" + structType + "}"))
			}
			return reflect.SliceOf(getType(typeName[2:]))
		}

		// Check for pointer types
		if len(typeName) > 1 && typeName[0] == '*' {
			return reflect.PointerTo(getType(typeName[1:]))
		}

		// Fallback to panic if unsupported type
		panic(fmt.Sprintf("unsupported field type: %s", typeName))
	}
}
