## Logs

```bash
go get github.com/zwk-app/go-tools/logs
```

```go
import (
	"github.com/zwk-app/go-tools/logs"
)

func test() {
  logs.SetLevelInfo()
  e := fmt.Errorf("some error")
  logs.Debug("title", "message not displayed (debug)", e)
}
```

## Tools

```bash
go get github.com/zwk-app/go-tools/tools
```

```go
import (
	"github.com/zwk-app/go-tools/tools"
)

func test() {
  empty := ""
  value := tools.Fallback(empty, "fallback value")
}
```
