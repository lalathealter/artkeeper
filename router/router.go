package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lalathealter/artkeeper/controllers"
)

const (
	apiurls        = "/api/urls"
	apicollections = "/api/collections"
)

func Use() *router {
	rt := &router{	}

	rt.setroute("/", "GET", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})


	rt.setroute(apiurls, "POST", controllers.PostURLhandler)
	rt.setroute(apiurls, "GET", controllers.GetURLHandler)
	rt.setroute(apiurls, "DELETE", controllers.DeleteURLHandler)

	apiurlslatest := fmt.Sprintf("%v/%v", apiurls, "latest")
	rt.setroute(apiurlslatest, "GET", controllers.GetLatestURLsHandler)

	rt.setroute(apicollections, "GET", controllers.GetCollectionHandler)
	rt.setroute(apicollections, "POST", controllers.PostCollectionHandler)
	rt.setroute(apicollections, "PUT", controllers.PutInCollectionHandler)
	rt.setroute(apicollections, "DELETE", controllers.DeleteCollectionHandler)
	return rt
}

type router map[int][]routeEntry

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedHandler := http.NotFound
	reqPathTokens := strings.Split(r.URL.Path, "/") 
	for _, rentry := range (*rt)[len(reqPathTokens)] {
		if rentry.doesMatchMethod(r) {
			score := rentry.doesMatchPath(reqPathTokens) 
			if score { 
				matchedHandler = rentry.Handler.ServeHTTP
			}
		}
	}

	matchedHandler(w, r)
}

func (rentry *routeEntry) doesMatchPath(requestPathTokens []string) (bool) {
	for i, val := range requestPathTokens {
		currToken := rentry.Path[i]
		if val != currToken && currToken != "*" {
			return false 
		}
	}
	return true
}

func (rt *router) setroute(p string, m string, hf http.HandlerFunc) {
	pathTokens := strings.Split(p, "/")
	rentry := routeEntry{
		Path:    pathTokens,
		Method:  strings.ToUpper(m),
		Handler: hf,
	}
	l := len(pathTokens)
	(*rt)[l] = append((*rt)[l], rentry)
}

type routeEntry struct {
	Path    []string
	Method  string
	Handler http.HandlerFunc
}

func (rentry *routeEntry) doesMatchMethod(r *http.Request) bool {
	return r.Method == rentry.Method 
}
