package http

import "net/http"

func Docs(prefix string) http.Handler {
	return http.StripPrefix(prefix, http.FileServer(http.Dir("./api")))
}
