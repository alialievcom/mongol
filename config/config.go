package config

import (
	"AliAlievMos/mongol/models"
	"fmt"
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
)

func LoadConfig(filename string) (*models.Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var cfg models.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling YAML: %w", err)
	}

	cfg.GeneratedStructMap = make(map[string]reflect.Type)

	for _, collection := range cfg.Mongo.Collections {
		details, exists := cfg.Mongo.Details[collection]
		if !exists {
			return nil, fmt.Errorf("%s structure not exist", details)
		}
		structType := generateStruct(details.Fields)
		cfg.GeneratedStructMap[collection] = structType
	}

	cfg.Api.SecretKey = []byte(cfg.Api.SecretKeyYML)

	return &cfg, nil
}

func generateStruct(fields []models.Field) reflect.Type {
	var structFields []reflect.StructField

	for _, field := range fields {
		if field.Type == "struct" {
			if field.Fields == nil {
				log.Panicf("fileds in %s is nil while type struct", field.Name)
			}
			refType := generateStruct(*field.Fields)
			structFields = append(structFields, reflect.StructField{
				Name: field.Name,
				Type: refType,
				Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(field.Name), strcase.ToSnake(field.Name), field.Tags)),
			})
			continue
		}
		structFields = append(structFields, reflect.StructField{
			Name: field.Name,
			Type: getType(field.Type),
			Tag:  reflect.StructTag(fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(field.Name), strcase.ToSnake(field.Name), field.Tags)),
		})
	}
	structFields = append(structFields, reflect.StructField{
		Name: "ID",
		Type: reflect.TypeOf(&primitive.ObjectID{}),
		Tag:  `json:"_id" bson:"_id"`,
	})

	return reflect.StructOf(structFields)
}
