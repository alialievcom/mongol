package config

import (
	"fmt"
	"github.com/AliAlievMos/mongol/models"
	"github.com/iancoleman/strcase"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	config, err := LoadConfig("./config_test.yml")
	if err != nil {
		t.Errorf("error loading config: %v", err)
	}

	expectedStructures := map[string][]models.Field{
		"publication": {
			{Name: "Title", Type: "string", Tags: "binding:\"required\""},
			{Name: "Id", Type: "int64", Tags: "binding:\"required\""},
		},
		"translations": {
			{Name: "Title", Type: "string", Tags: "binding:\"required\""},
			{Name: "Id", Type: "int64", Tags: "binding:\"required\""},
		},
		"authors": {
			{Name: "Title", Type: "string", Tags: "binding:\"required\""},
			{Name: "Id", Type: "int64", Tags: "binding:\"required\""},
		},
	}

	for collection, expectedFields := range expectedStructures {
		generatedStruct, ok := config.GeneratedStructMap[collection]
		if !ok {
			t.Errorf("expected structure for collection %s not found", collection)
			continue
		}

		generatedType := reflect.TypeOf(generatedStruct).Elem()

		if generatedType.NumField() != len(expectedFields) {
			t.Errorf("expected %d fields in structure for collection %s, got %d",
				len(expectedFields), collection, generatedType.NumField())
			continue
		}

		for i, expectedField := range expectedFields {
			field := generatedType.Field(i)
			if field.Name != expectedField.Name {
				t.Errorf("expected field name %s in structure for collection %s, got %s",
					expectedField.Name, collection, field.Name)
			}
			if field.Type.Name() != expectedField.Type {
				t.Errorf("expected field type %s for field %s in structure for collection %s, got %s",
					expectedField.Type, expectedField.Name, collection, field.Type.Name())
			}
			if string(field.Tag) != fmt.Sprintf(`json:"%s" bson:"%s" %s`, strcase.ToSnake(expectedField.Name), strcase.ToSnake(expectedField.Name), expectedField.Tags) {
				t.Errorf("expected field tag %s for field %s in structure for collection %s, got %s",
					fmt.Sprintf(`json:"%s" binding:"%s"`, strcase.ToSnake(expectedField.Name), expectedField.Tags),
					expectedField.Name, collection, string(field.Tag))
			}
		}
	}
}
