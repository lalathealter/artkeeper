package router

import (
	"net/http"
	"strings"

	"github.com/lalathealter/artkeeper/controllers"
)

func Use() *router {
	rt := &router{}
	rt.setroute("/", "GET", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})
	rt.setroute("/api/urls", "POST", controllers.PostURLhandler)
	rt.setroute("/api/urls", "GET", controllers.GetOneURLHandler)
	return rt
}

type router struct {
	routes []routeEntry
}

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, rentry := range rt.routes {
		if rentry.doesmatch(r) {
			rentry.Handler.ServeHTTP(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func (rt *router) setroute(p string, m string, hf http.HandlerFunc) {
	rentry := routeEntry{
		Path:    p,
		Method:  strings.ToUpper(m),
		Handler: hf,
	}
	rt.routes = append(rt.routes, rentry)
}

type routeEntry struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

func (rentry *routeEntry) doesmatch(r *http.Request) bool {
	return r.Method == rentry.Method && r.URL.Path == rentry.Path
}
