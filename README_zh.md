# SwaggerAPI

`SwaggerAPI` 是一个简单的 API，允许您为API创建SwaggerUI，为项目开发提供帮助。它基于[Swagger 2.0规范]()。
通过读取进程中的grpc服务，自动生成swagger文档数据，同时提供swagger-ui服务，方便查看文档。

推荐使用[gkit](https://github.com/ml444/gkit)框架，可以自动生成swagger文档，同时提供swagger-ui服务。

本项目是[gkit](https://github.com/ml444/gkit)的需求衍生品，也可以单独使用。

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
	//userpb.RegisterServerWithHTTP(srv)
	
	if err = srv.Start(context.Background()); err != nil {
		log.Error(err.Error())
		return
	}
}
```
```shell
$ go run main.go
```
**Visit**：http://localhost:5050/swagger
