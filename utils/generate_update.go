package utils

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

var (
	structType = reflect.ValueOf(struct{}{}).Type()
)

func GenerateUpdateBson(v interface{}) (bson.M, error) {
	update := bson.M{}
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", val.Kind())
	}

	t := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := t.Field(i)

		if typeField.PkgPath != "" {
			continue
		}

		fieldName := typeField.Tag.Get("bson")
		if fieldName == "" {
			fieldName = strings.ToLower(typeField.Name)
		}
		if fieldName == "_id" {
			continue
		}

		if field.Kind() == reflect.Ptr && field.IsNil() {
			update[getFieldBsonName(fieldName, typeField)] = nil
			continue
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}

		if field.Type() == structType {
			embeddedDoc, err := GenerateUpdateBson(field.Interface())
			if err != nil {
				return nil, err
			}
			for k, v := range embeddedDoc {
				update[getFieldBsonName(fieldName+"."+k, typeField)] = v
			}
		} else {
			update[getFieldBsonName(fieldName, typeField)] = field.Interface()
		}
	}
	updateFinal := bson.M{"$set": update}
	return updateFinal, nil
}

func GenerateUpdateBsonInterface(data interface{}) (bson.M, error) {
	update := bson.M{}
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Map {
		return nil, fmt.Errorf("expected a map, got %s", val.Kind())
	}

	iter := val.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		if k.Kind() == reflect.String {
			fieldName := k.String()
			update[fieldName] = v.Interface()
		}
	}

	return bson.M{"$set": update}, nil
}

func getFieldBsonName(fieldName string, f reflect.StructField) string {
	tag := f.Tag.Get("bson")
	if tag == "" {
		return fieldName
	}
	tagParts := strings.Split(tag, ",")
	if len(tagParts) > 0 && tagParts[0] != "" {
		return tagParts[0]
	}
	return fieldName
}
