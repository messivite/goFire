package handlers

import "net/http"

var Registry = make(map[string]http.HandlerFunc)

func Register(name string, fn http.HandlerFunc) {
	Registry[name] = fn
}

func Get(name string) http.HandlerFunc {
	return Registry[name]
}
