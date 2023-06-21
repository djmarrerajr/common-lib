## My Sandbox Repo
### package: `app`
---
<br>

### Here are some examples to get you started:
<br>

```go
 // A simple HTTP Api Server
 //--------------------------------------------------------------------------------
 app.WithNewHttpsServer("0.0.0.0", "8080",
 	api.WithRouteHandler("/", indexHandler),
 )
```

```go
// A simple HTTPS Api Server
//--------------------------------------------------------------------------------
 app.WithNewHttpsServer("0.0.0.0", "8443",
 	api.WithRouteHandler("/", indexHandler),
 )
```

```go
// A simple HTTPS Api Server w/mTLS
//--------------------------------------------------------------------------------
 app.WithNewHttpsServer("0.0.0.0", "8443", "./certs/localhost.crt", "./certs/localhost.key",
 	api.WithMtlsEnforcedCaCert("./certs/ca.crt"),
 	api.WithRouteHandler("/", indexHandler),
 )
```