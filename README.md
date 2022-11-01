# golog

Simple golang rotated log.

### Install

```go 
go get github.com/black40x/golog
```

### Example

```go 
package main

import (
	"golog"
)

func main() {
	rl := golog.NewLogger(&golog.Options{
		MaxSize: 10 << 20,
		LogName: "my-log",
		Daily:   true,
	}, golog.Ldate|golog.Ltime)
	defer rl.Close()

	rl.Println("My log")
	rl.Info("My tagged log Info")
	rl.Error("My tagged log Error")
	rl.Warning("My tagged log Warning")
}
```
