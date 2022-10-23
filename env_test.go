package env

import (
	"fmt"
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
		"with nested value and default and replace": {[]byte(`KEY=VALUE
KEY2=${KEY3:=DEFAULT}`), map[string]string{"KEY": "VALUE", "KEY2": "DEFAULT", "KEY3": "DEFAULT"}},
	}

	for name, test := range data {
		t.Run(name, func(t *testing.T) {
			env := &Env{Kv: make(map[string]string), Keys: make([]string, 0)}
			env.Parser(test.input)
			fmt.Println(env)
			for key, value := range test.expected {
				if env.Kv[key] != value {
					t.Errorf("\nExpected -> %s\nGot -> %s", value, env.Kv[key])
				}
			}
		})
	}
}
