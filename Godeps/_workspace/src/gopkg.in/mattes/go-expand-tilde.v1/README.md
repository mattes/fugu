# Expand ~ in path

```
go get gopkg.in/mattes/go-expand-tilde.v1
```

```go
import (
  "gopkg.in/mattes/go-expand-tilde.v1"
)

func main() {
  path, err := tilde.Expand("~/path/to/whatever")
}
```