package config

import (
	"fmt"
	"reflect"
)

type Config struct {
	Port     string `env:"SERVER_PORT" envDocs:"The port which the service will listen to" envDefault:"8080"`
	BasePath string `env:"SERVER_BASE_PATH" envDocs:"The base path of this service" envDefault:"/player"`

	MongoURL            string `env:"MONGO_URL" envDocs:"mongo url" envDefault:"mongodb://localhost:27017"`
	MongoUserName       string `env:"MONGO_USER_NAME" envDocs:"mongo user name"`
	MongoPassword       string `env:"MONGO_PASSWORD" envDocs:"mongo password"`
	MongoDatabase       string `env:"MONGO_DB" envDocs:"Mongo Database name" envDefault:"player"`
	MongoContextTimeout int    `env:"CONTEXT_TIMEOUT_DURATION" envDocs:"context timeout duration, unit second" envDefault:"10"`
}

func (envVar Config) HelpDocs() []string {
	reflectEnvVar := reflect.TypeOf(envVar)
	doc := make([]string, 1+reflectEnvVar.NumField())
	doc[0] = "Environment variables config:"
	for i := 0; i < reflectEnvVar.NumField(); i++ {
		field := reflectEnvVar.Field(i)
		envName := field.Tag.Get("env")
		envDefault := field.Tag.Get("envDefault")
		envDocs := field.Tag.Get("envDocs")
		doc[i+1] = fmt.Sprintf("  %v\t %v (default: %v)", envName, envDocs, envDefault)
	}
	return doc
}
