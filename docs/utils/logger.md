## `logger`
### structured logger
---
<br>

Example usage:
```go
1: func main() {
2: 	 var logger = utils.NewLogger("DEBUG").Named("my-app")
3:	 var userId = "bruno"
4:
5:	 logger.Debugw("login attempt", "user", userId)
6: }
```
The above example creates a new, DEBUG-level, structured logger named `my-app` the result of line 5 would be:

```console
{"level":"debug","ts":1683077739.5557911,"logger":"my-app","caller":"cmd/main.go:9","msg":"login attempt","user":"bruno"}
```

Once a logger has been created it is possible to enhance its usage by, among other things:
- creating `child` loggers (e.g. `logger.Named("sub-module")`)
- adding contextual information via the `.WithCtx()` method (i.e. application version information, etc.)