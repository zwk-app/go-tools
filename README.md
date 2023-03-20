## Logs

```bash
go get github.com/zwk-app/go-tools/logs
```

```go
package main

import (
	"github.com/zwk-app/go-tools/logs"
)

func main() {
  logs.SetLevelInfo()
  e := fmt.Errorf("some error in log and StdErr")
  logs.Debug("title", "message not displayed (debug but Info log level)", e)
}
```

## Tools

```bash
go get github.com/zwk-app/go-tools/tools
```

```go
package main

import (
	"github.com/zwk-app/go-tools/tools"
)

func main() {
  empty := ""
  value := tools.Fallback(empty, "fallback value")
}
```
