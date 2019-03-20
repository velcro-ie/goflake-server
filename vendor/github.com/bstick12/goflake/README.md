# goflake

[![CircleCI](https://img.shields.io/circleci/project/github/bstick12/goflake.svg)](https://circleci.com/gh/bstick12/goflake) [![Codecov](https://img.shields.io/codecov/c/github/bstick12/goflake.svg)](https://codecov.io/gh/bstick12/goflake)

A flake id generator based on the [TimebasedUUIDGenerator](https://github.com/elastic/elasticsearch/blob/master/core/src/main/java/org/elasticsearch/common/TimeBasedUUIDGenerator.java) in [Elasticsearch](https://github.com/elastic/elasticsearch)

## Sample Usage
```
package main

import (
	. "github.com/bstick12/goflake"
	"fmt"
)

func main() {
	generator := GetGoFlakeInstance()
    uuid := generator.GetBase64UUID()
    fmt.Println(uuid)
}

```