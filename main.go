package main

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Kv   map[string]string
	Keys []string
}

func (env *Env) String() string {
	var str string
	for key, value := range env.Kv {
		str += fmt.Sprintf("%s=%s\n", key, value)
	}
	return str
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

		parts := strings.SplitN(line, "=", 2)
		key, value := strings.TrimSpace(parts[0]), strings.Trim(parts[1], `" `)
		env.Kv[key] = value
	}

	for key, value := range env.Kv {
		env.Kv[key] = env.SolveNestedValues(value)
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
			var replacer string

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

			if hasDefault {
				replacer = defaultValue
			} else {
				var ok bool
				replacer, ok = env.Kv[key]
				if !ok {
					panic(fmt.Sprintf("Key %s not found in `%s`", key, value))
				}
			}

			if replace {
				env.Kv[key] = replacer
			}

			target := "${" + parameter + "}"
			value = strings.Replace(value, target, replacer, 1)
		}
	}
	return value
}

func Load() {
	data, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println(err)
	}
	env := &Env{Kv: make(map[string]string)}
	env.Parser(data)
	fmt.Println("Loading environment values from .env file:")
	fmt.Println(env)
	for key, value := range env.Kv {
		err := os.Setenv(key, value)
		if err != nil {
			panic(err)
		}
	}
}

func init() {
	Load()
}

func main() {
	fmt.Println("In main program:")
}
