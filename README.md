# go.env

## Features

- [x] load `.env` file
- [x] support string syntax
- [x] nested variables using `${}` syntax
- [x] set default value using `:=` and `:-`

## Usage

```go
package main

import (
  "github.com/ekil1100/go.env"
)

func init() {
  env.Load()
}

func main() {
}
```

## Reference

- https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
