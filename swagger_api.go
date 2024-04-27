package swaggerapi

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/ml444/swaggerapi/generator"
)

//go:embed dist/*
var dist embed.FS

type Router interface {
	HandlePrefix(prefix string, h http.Handler)
	HandleFunc(path string, h http.HandlerFunc)
}

func RegisterAPI(r Router, opts ...generator.Option) {
	ss := NewServer(opts...)
	r.HandlePrefix("/swagger/", NewSwaggerHandler())
	r.HandlePrefix("/swagger-query/service/{name}", http.Handler(http.HandlerFunc(ss.GetServiceDesc)))
	r.HandleFunc("/swagger-query/services", ss.ListServices)
}

func NewSwaggerHandler() http.Handler {
	fSys, err := fs.Sub(dist, "dist")
	if err != nil {
		println(err.Error())
	}
	handler := http.FileServer(http.FS(fSys))
	return http.StripPrefix("/swagger/", handler)
}
