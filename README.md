## Proprietary Tenders - Gift Cards
### Common Library
---
<br>

Table of Contents:

* [Overview]()
* [Module Guides](docs/README.md)
	* [errs](docs/errs/README.md)
	* [utils](docs/utils/README.md)



<br>

---

#### A very simple, albeit contrived, example:
<br>

```go
package main

import (
	"github.com/djmarrerajr/common-lib/errs"
	"github.com/djmarrerajr/common-lib/utils"
)

func main() {
	logger := utils.NewLogger("INFO")

	err := validate()
	if err != nil {
		logger.Error("startup failed", err)
	}
}

func validate() error {
	return errs.New(errs.ErrTypeValidation, "validation failed")
}
```

<br>

... and the output:
```shell
$ go run main.go 2>&1 | jq

{
  "level": "error",
  "ts": 1683114890.155433,
  "caller": "main.go:13",
  "msg": "startup failed",
  "error.message": "validation failed",
  "stacktrace": "main.main\n\t/main.go:13\nruntime.main\n\t/usr/local/go/src/runtime/proc.go:250"
}

```