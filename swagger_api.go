package swaggerapi

import (
	"embed"
	"net/http"

	"github.com/ml444/swaggerapi/generator"
)

//go:embed dist/*
var dist embed.FS

type Router interface {
	HandlePrefix(method, prefix string, h http.Handler)
	HandleFunc(method, path string, h http.HandlerFunc)
}

func RegisterApi(r Router, opts ...generator.Option) {
	ss := NewServer(opts...)
	r.HandlePrefix(http.MethodGet, "/swagger/", NewSwaggerHandler())
	r.HandlePrefix(http.MethodGet, "/query/service/{name}", http.Handler(http.HandlerFunc(ss.GetServiceDesc)))
	r.HandleFunc(http.MethodGet, "/query/services", ss.ListServices)
}

func NewSwaggerHandler() http.Handler {
	fs := http.FileServer(http.FS(dist))
	return http.StripPrefix("/swagger/", fs)
}
