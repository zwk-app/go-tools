## Logs

```bash
go get github.com/zwk-app/zwk-tools/logs
```

```bash
go get github.com/zwk-app/zwk-tools/tools
```

```bash
go get github.com/zwk-app/zwk-tools/timer
```

```bash
go get github.com/zwk-app/zwk-tools/tests
```

```bash
go get github.com/zwk-app/zwk-tools
```

```go
package main

import (
	"github.com/zwk-app/zwk-tools/logs"
)

func main() {
  logs.SetLevelInfo()
  e := fmt.Errorf("some error in log and StdErr")
  logs.Debug("title", "message not displayed (debug but Info log level)", e)
}
```
