quickproxy
==========

## Examples

### Record request/response elapsed time.

```go
proxy server
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
