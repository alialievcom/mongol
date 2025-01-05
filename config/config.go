package config

import (
	"fmt"
	"github.com/AliAlievMos/mongol/models"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
)

// LoadConfig loads the configuration from the YAML file and generates structs based on the fields
func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	cfg.GeneratedStructMap = make(map[string]models.Collection)

	for _, collection := range cfg.Mongo.Collections {
		details, exists := cfg.Mongo.Details[collection]
		if !exists {
			return nil, fmt.Errorf("%v structure not exist", details)
		}
		structType := generateStructWrapper(details.Fields)
		cfg.GeneratedStructMap[collection] = models.Collection{
			Model:        structType,
			SortBy:       details.SortBy,
			QueryFilters: details.GetQueryFilters(),
		}
	}

	cfg.Api.SecretKey = []byte(cfg.Api.SecretKeyYML)

	return &cfg, nil
}

// generateStructWrapper generates a struct wrapper with an additional ID field
func generateStructWrapper(fields []models.Field) reflect.Type {
	structFields := generateStruct(fields)
	structFields = append(structFields, reflect.StructField{
		Name: "ID",
		Type: reflect.TypeOf(&primitive.ObjectID{}),
		Tag:  `json:"_id" bson:"_id"`,
	})

	return reflect.StructOf(structFields)
}

// generateStruct generates a list of struct fields from the given fields
func generateStruct(fields []models.Field) []reflect.StructField {
	var structFields []reflect.StructField

	for _, field := range fields {
		if field.Type == "struct" || field.Type == "*struct" {
			if field.Fields == nil && field.Type == "struct" {
				log.Panicf("fields in %s are nil while type is struct", field.Name)
			}
			structField := generateStruct(*field.Fields)
			refType := reflect.StructOf(structField)
			if field.Type == "*struct" {
				refType = reflect.PointerTo(refType)
			}
			structFields = append(structFields, reflect.StructField{
				Name: field.Name,
				Type: refType,
				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(field.Name), strcase.ToSnake(field.Name), field.Tags)),
			})
			continue
		}
		// Handle slices like []struct
		if field.Type == "[]struct" {
			// Recursively handle the inner struct fields
			if field.Fields == nil {
				log.Panicf("fields in %s are nil while type is []struct", field.Name)
			}
			structField := generateStruct(*field.Fields)
			refType := reflect.SliceOf(reflect.StructOf(structField))
			structFields = append(structFields, reflect.StructField{
				Name: field.Name,
				Type: refType,
				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(field.Name), strcase.ToSnake(field.Name), field.Tags)),
			})
			continue
		}

		// Handle other types like string, int, map, etc.
		structFields = append(structFields, reflect.StructField{
			Name: field.Name,
			Type: getType(field.Type),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(field.Name), strcase.ToSnake(field.Name), field.Tags)),
		})
	}
	return structFields
}
