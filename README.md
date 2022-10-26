# goenv

Import environment variables from `.env` file.

## Features

- [x] no dependency
- [x] load `.env` file
- [x] support string syntax
- [x] support comment
- [x] support inline comment
- [x] nested variables using `${}` syntax
- [x] set default value using `:=` and `:-`
- [x] Don't load the `.env` if it's not present
- [ ] Support custom path to the `.env` file
- [x] load `.env.xxx` depends on `APP_ENV` value

## Install

```shell
go get github.com/ekil1100/goenv
```

## Usage

Assume you have a `.env` file in the root of your project:

```shell
DB_USER=postgres
DB_PASS=xyz
DB_HOST=localhost
DB_PORT=5432
DB_NAME=db
DATABASE_URL="postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?pool_timeout=30&connection_limit=22"
```

Then in your can add `env.Load()` in your entry file:

```go
package main

import (
  "fmt"
  "log"
  "os"

  env "github.com/ekil1100/goenv"
)

func init() {
  err := env.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }
}

func main() {
  fmt.Println(os.Getenv("DATABASE_URL")) // postgresql://postgres:xyz@localhost:5432/db?pool_timeout=30&connection_limit=22
}
```

### Nested variables

You can use `${}` to reference value.

### Default value

You can use `${KEY:defaultValue}` to set default.

`${parameter:-value}` is same as `${parameter:value}`, they both just use `word` to substitute `parameter`
if `parameter` is unset or null.

`${parameter:-value}` on the other hand will create a new value for you if `parameter` is unset or null.

## Reference

- https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
