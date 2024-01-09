# SwaggerAPI
`SwaggerAPI` is a simple API that allows you to create SwaggerUI for the API to help project development. 
It is based on[Swagger 2.0 Specification]()。

By reading the grpc service in the process, swagger document data is automatically generated, 
and the swagger-ui service is provided to facilitate document viewing.

> It is recommended to use the [gkit](https://github.com/ml444/gkit) framework, which can automatically generate swagger documents and provide swagger-ui services.

This project is a demand derivative of [gkit](https://github.com/ml444/gkit) and can also be used independently.

## Installation
```bash
$ go get github.com/ml444/swaggerapi
```
## Usage
```go
package main

import (
    "github.com/ml444/gkit/log"
    "github.com/ml444/gkit/transport/httpx"
    "github.com/ml444/swaggerapi"
)

func main() {
	var err error
	srv := httpx.NewServer(
		httpx.Timeout(time.Duration(hcfg.Timeout)*time.Millisecond),
		httpx.Address(hcfg.HTTPAddr),
	)
	// Create a new SwaggerAPI instance
	swaggerapi.RegisterAPI(srv)
	
	// Register your API handlers
	userpb.RegisterServerWithHTTP(srv)
	
	if err = srv.Start(context.Background()); err != nil {
		log.Error(err.Error())
		return
	}
}
```

```shell
$ go run main.go
```
访问：http://localhost:5050/swagger