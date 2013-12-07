Go Library for a easy HTTP proxy
==========

Go Library for a easy HTTP proxy

## Examples

### Record request/response elapsed time.

```go
package main

import (
	"fmt"
	"github.com/reoring/quickproxy"
)

func main() {
	quickproxy.Prepare(map[string]string{"port": "8081"})

	quickproxy.OnDone(func(doneRequestData *quickproxy.DoneRequestData) {
		fmt.Print(doneRequestData.Request.URL)
		fmt.Print(": ")
		fmt.Print(doneRequestData.RoundTripTime.ElapsedTime())
		fmt.Print("\n")
	})

	quickproxy.Run()
}
```
