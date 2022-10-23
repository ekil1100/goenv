package main

import (
	"fmt"
	"os"
	"strings"
)

type Env struct {
	Envs map[string]string
}

func (env *Env) String() string {
	var str string
	for key, value := range env.Envs {
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
		env.Envs[key] = value
	}

	for key, value := range env.Envs {
		env.Envs[key] = env.SolveNestedValues(value)
	}
}

func (env *Env) SolveNestedValues(value string) string {
	for i := 0; i < len(value); i++ {
		if value[i] == '$' {
			end := strings.Index(value[i:], "}") + i
			key := value[i+2 : end]
			if value, ok := env.Envs[key]; !ok {
				panic(fmt.Sprintf("Key %s not found from `%s`", key, value))
			}
			value = strings.Replace(value, fmt.Sprintf("${%s}", key), env.Envs[key], 1)
		}
	}
	return value
}

func Load() {
	data, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println(err)
	}
	env := &Env{Envs: make(map[string]string)}
	env.Parser(data)
	fmt.Println("Loading environment values from .env file:")
	fmt.Println(env)
	for key, value := range env.Envs {
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
