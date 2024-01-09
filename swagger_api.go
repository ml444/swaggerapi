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
	HandlePrefix(method, prefix string, h http.Handler)
	HandleFunc(method, path string, h http.HandlerFunc)
}

func RegisterAPI(r Router, opts ...generator.Option) {
	ss := NewServer(opts...)
	r.HandlePrefix(http.MethodGet, "/swagger/", NewSwaggerHandler())
	r.HandlePrefix(http.MethodGet, "/swagger-query/service/{name}", http.Handler(http.HandlerFunc(ss.GetServiceDesc)))
	r.HandleFunc(http.MethodGet, "/swagger-query/services", ss.ListServices)
}

func NewSwaggerHandler() http.Handler {
	fSys, err := fs.Sub(dist, "dist")
	if err != nil {
		println(err.Error())
	}
	handler := http.FileServer(http.FS(fSys))
	return http.StripPrefix("/swagger/", handler)
}
