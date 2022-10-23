package env

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Kv   map[string]string // key value store
	Keys []string          // to keep the order of the keys
}

func Filter(list []string, fn func(string) bool) (result []string) {
	for _, item := range list {
		if fn(item) {
			result = append(result, item)
		}
	}
	return result
}

func (env *Env) String() string {
	var str string
	for _, key := range env.Keys {
		value := env.Kv[key]
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	return str
}

func (env *Env) Add(key, value string) {
	env.Keys = Filter(env.Keys, func(k string) bool {
		return k != key
	})
	env.Keys = append(env.Keys, key)
	env.Kv[key] = value
}

func (env *Env) Parser(data []byte) {
	for lineNumber, line := range strings.Split(string(data), "\n") {
		// trim spaces
		line = strings.TrimSpace(line)

		// skip # comments
		if strings.HasPrefix(line, "#") {
			continue
		}

		// skip empty lines
		if line == "" {
			continue
		}

		// panic if no = sign, invalid syntax
		if !strings.Contains(line, "=") {
			panic(fmt.Sprintf("Line %d: `%s` is not a valid environment variable\n", lineNumber+1, line))
		}

		// split key and value and store into a map
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Split(parts[1], "#")[0]
		value = strings.Trim(value, "\" ")
		env.Add(key, value)
	}

	// solve nested values
	for _, key := range env.Keys {
		value := env.Kv[key]
		env.Add(key, env.SolveNestedValues(value))
	}
}

func (env *Env) SolveNestedValues(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] == '$' {
			end := strings.Index(value[i:], "}") + i
			parameter := value[i+2 : end]
			spitedParameter := strings.Split(parameter, ":")

			var key, defaultValue string
			var replace, hasDefault bool

			if len(spitedParameter) > 1 {
				hasDefault = true
				key = spitedParameter[0]
				if strings.HasPrefix(spitedParameter[1], "=") {
					replace = true
					defaultValue = spitedParameter[1][1:]
				} else if strings.HasPrefix(spitedParameter[1], "-") {
					defaultValue = spitedParameter[1][1:]
				} else {
					defaultValue = spitedParameter[1]
				}
			} else {
				key = spitedParameter[0]
			}

			replacer, ok := env.Kv[key]
			if !ok {
				if hasDefault {
					replacer = defaultValue
				} else {
					panic(fmt.Sprintf("Key %s not found in `%s`", key, value))
				}
			}

			if replace {
				env.Add(key, replacer)
			}

			target := "${" + parameter + "}"
			value = strings.Replace(value, target, replacer, 1)
		}
	}
	return value
}

func Load() error {
	data, err := os.ReadFile(".env")
	if err != nil {
		return err
	}
	env := &Env{Kv: make(map[string]string), Keys: make([]string, 0)}
	env.Parser(data)
	fmt.Println("Loading environment values from .env file:")
	fmt.Println(env)
	for key, value := range env.Kv {
		err := os.Setenv(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}
