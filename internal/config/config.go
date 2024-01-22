package config

import (
	"fmt"
	"os"
	"reflect"
)

type Config struct {
	ServerAddr        string `env:"SERVER_ADDR"`
	ServerRoutePrefix string `env:"SERVER_ROUTE_PREFIX"`
	VaultAddr         string `env:"VAULT_ADDR"`
	VaultEngine       string `env:"VAULT_ENGINE"`
	VaultToken        string `env:"VAULT_TOKEN"`
}

func (c *Config) Load() error {
	v := reflect.ValueOf(c).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)

		key, ok := f.Tag.Lookup("env")
		if !ok {
			return fmt.Errorf("field '%s' does not have env tag", f.Name)
		}

		env, ok := os.LookupEnv(key)
		if !ok {
			return fmt.Errorf("environment variable '%s' does not exist", key)
		}

		s := v.Field(i)
		if !s.CanSet() {
			return fmt.Errorf("field '%s' is not settable", f.Name)
		}

		if s.Kind() != reflect.String {
			return fmt.Errorf("field '%s' kind is not string", f.Name)
		}

		s.SetString(env)
	}

	return nil
}
