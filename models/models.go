package models

import (
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type Field struct {
	Name     string   `yaml:"Name"`
	Type     string   `yaml:"Type"`
	Tags     string   `yaml:"Tags"`
	IsHidden bool     `yaml:"IsHidden"`
	Fields   *[]Field `yaml:"fields"`
}

type CollectionConfig struct {
	Fields            []Field           `yaml:"fields"`
	SortBy            string            `yaml:"sort_by"`
	NecessaryAuthRole map[string]string `yaml:"necessary_auth_role"`
}

func (cc *CollectionConfig) GetQueryFilters() []string {
	var filters = make([]string, len(cc.Fields))
	for i, field := range cc.Fields {
		filters[i] = strcase.ToSnake(field.Name)
	}
	return filters
}

type MongoConfig struct {
	DB          string `yaml:"db"`
	URI         string
	Collections []string                    `yaml:"collections"`
	Details     map[string]CollectionConfig `yaml:",inline"`
	Users       []User                      `yaml:"users"`
	Auth        Auth                        `yaml:"auth"`
}

type Api struct {
	Port         string `yaml:"port"`
	SecretKeyYML string `yaml:"secret_key"`
	SecretKey    []byte `yaml:"-"`
	Origin       string `yaml:"origin"`
	Headers      string `yaml:"headers"`
	Methods      string `yaml:"methods"`
}
type Auth struct {
	AuthCollection *string `yaml:"auth_collection"`
	AuthLocation   string  `yaml:"auth_location"`
}
type Config struct {
	Name               string                `yaml:"name"`
	Api                Api                   `yaml:"api"`
	Mongo              MongoConfig           `yaml:"mongo"`
	GeneratedStructMap map[string]Collection `yaml:"-"`
}

type Collection struct {
	Model             reflect.Type
	SortBy            string
	NecessaryAuthRole map[string]string
	QueryFilters      []string
}

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Login    string             `yaml:"login" bson:"login" json:"login"`
	Roles    []string           `yaml:"roles" bson:"roles" json:"roles"`
	Password string             `yaml:"password" bson:"password"`
}

type Token struct {
	Access string `json:"access"`
	User
}
