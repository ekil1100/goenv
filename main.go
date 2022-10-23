package main

import (
	"fmt"
	"os"
	"strings"
)

func Load() {
	data, err := os.ReadFile(".env")
	if err != nil {
		fmt.Println(err)
	}
	parser(data)
}

func parser(data []byte) {
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, "=") {
			parts := strings.Split(line, "=")
			err := os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			if err != nil {
				panic(err)
			}
		}
	}
}

func init() {
	Load()
}

func main() {
	fmt.Println("In main program:")
	fmt.Println(os.Environ()[0])
}
