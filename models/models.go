package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type Field struct {
	Name   string   `yaml:"Name"`
	Type   string   `yaml:"Type"`
	Tags   string   `yaml:"Tags"`
	Fields *[]Field `yaml:"fields"`
}

type CollectionConfig struct {
	Fields []Field `yaml:"fields"`
}

type MongoConfig struct {
	DB          string                      `yaml:"db"`
	URI         string                      `yaml:"uri"`
	Collections []string                    `yaml:"collections"`
	Details     map[string]CollectionConfig `yaml:",inline"`
	Users       []User                      `yaml:"users"`
}

type Api struct {
	Port         string `yaml:"port"`
	SecretKeyYML string `yaml:"secret_key"`
	SecretKey    []byte `yaml:"-"`
}

type Config struct {
	Name               string                  `yaml:"name"`
	Api                Api                     `yaml:"api"`
	Mongo              MongoConfig             `yaml:"mongo"`
	GeneratedStructMap map[string]reflect.Type `yaml:"-"`
}

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"_id"`
	Login    string             `yaml:"login" bson:"login"`
	Password string             `yaml:"password" bson:"password"`
}

type Token struct {
	Access string `json:"access"`
}
