package env

import (
	"os"
	"testing"
)

func TestEnv_Parser(t *testing.T) {
	data := map[string]struct {
		input    []byte
		expected map[string]string
	}{
		"simple":                          {[]byte("KEY=VALUE"), map[string]string{"KEY": "VALUE"}},
		"with space":                      {[]byte("KEY = VALUE"), map[string]string{"KEY": "VALUE"}},
		"with space and quotes":           {[]byte("KEY = \"VALUE\""), map[string]string{"KEY": "VALUE"}},
		"with space and quotes and space": {[]byte("KEY = \"VALUE\" "), map[string]string{"KEY": "VALUE"}},
		"with inline comment":             {[]byte("KEY = \"VALUE\" # comment"), map[string]string{"KEY": "VALUE"}},
		"with comment": {[]byte(`KEY=VALUE
# COMMENT`), map[string]string{"KEY": "VALUE"}},
		"with nested value": {[]byte(`KEY=VALUE
KEY2=${KEY}`), map[string]string{"KEY": "VALUE", "KEY2": "VALUE"}},
		"with nested value and default": {[]byte(`KEY=VALUE
KEY2=${KEY:-DEFAULT}`), map[string]string{"KEY": "VALUE", "KEY2": "VALUE"}},
		"with nested value and default 2": {[]byte(`KEY=VALUE
KEY2=${KEY3:-DEFAULT}`), map[string]string{"KEY": "VALUE", "KEY2": "DEFAULT"}},
		"with nested value and default 3": {[]byte(`KEY=VALUE
KEY2=${KEY3:DEFAULT}`), map[string]string{"KEY": "VALUE", "KEY2": "DEFAULT"}},
		"with nested value and default and replace": {[]byte(`KEY=VALUE
KEY2=${KEY3:=DEFAULT}`), map[string]string{"KEY": "VALUE", "KEY2": "DEFAULT", "KEY3": "DEFAULT"}},
		"with nested value 2": {[]byte(`DB_USER=postgres
DB_PASS=xyz
DB_HOST=localhost
DB_PORT=5432
DB_NAME=db
DATABASE_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?pool_timeout=30&connection_limit=22"`), map[string]string{"DB_USER": "postgres", "DB_PASS": "xyz", "DB_HOST": "localhost", "DB_PORT": "5432", "DB_NAME": "db", "DATABASE_URL": "postgresql://postgres:xyz@localhost:5432/db?pool_timeout=30&connection_limit=22"}},
		"default with null": {[]byte(`KEY=
KEY2=${KEY:-DEFAULT}`), map[string]string{"KEY": "", "KEY2": "DEFAULT"}},
	}

	for name, test := range data {
		t.Run(name, func(t *testing.T) {
			env := &Env{Kv: make(map[string]string), Keys: make([]string, 0)}
			env.Parser(test.input)
			for key, value := range test.expected {
				if env.Kv[key] != value {
					t.Errorf("\nExpected -> %s\nGot -> %s", value, env.Kv[key])
				}
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		appEnv   string
		expected map[string]string
	}{
		{name: "simple", expected: map[string]string{"ENV": "default"}},
		{name: "with app env", appEnv: "local", expected: map[string]string{"ENV": "local"}},
		{name: "with filename", filename: ".env.local", expected: map[string]string{"ENV": "local"}},
		{name: "with app env and filename", appEnv: "local", filename: ".env.prod", expected: map[string]string{"ENV": "prod"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_ENV", tt.appEnv)
			var err error
			if tt.filename != "" {
				err = Load(tt.filename)
			} else {
				err = Load()
			}
			if err != nil {
				t.Errorf("Load() error = %v", err)
			}
			for k, v := range tt.expected {
				if os.Getenv(k) != v {
					t.Errorf("\nExpected -> %s=%s\nGot -> %s=%s", k, v, k, os.Getenv(k))
				}
			}
		})
	}
}
