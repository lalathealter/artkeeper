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

	apiurlsone := appendPath(apiurls, "*") 
	rt.setroute(apiurlsone, "POST", controllers.PostURLhandler)
	rt.setroute(apiurlsone, "GET", controllers.GetURLHandler)
	rt.setroute(apiurlsone, "DELETE", controllers.DeleteURLHandler)

	apiurlslatest := appendPath(apiurls, "latest") 
	rt.setroute(apiurlslatest, "GET", controllers.GetLatestURLsHandler)

	apicollectionsone := appendPath(apicollections, "*")
	rt.setroute(apicollectionsone, "GET", controllers.GetCollectionHandler)
	rt.setroute(apicollectionsone, "POST", controllers.PostCollectionHandler)
	rt.setroute(apicollectionsone, "PUT", controllers.PutInCollectionHandler)
	rt.setroute(apicollectionsone, "DELETE", controllers.DeleteCollectionHandler)
	apicollectionsurlsone := appendPath(apicollectionsone, "urls/*")
	rt.setroute(apicollectionsurlsone, "DELETE", controllers.DeleteURLFromCollection)
	fmt.Println(*rt)
	return rt
}

type router map[int][]routeEntry

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedHandler := http.NotFound
	reqPathTokens := parsePathTokens(r.URL.Path)
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
	pathTokens := parsePathTokens(p) 
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

func parsePathTokens(path string) []string {
	return strings.Split(path, "/")
}


func appendPath(base, next string) string {
	if next[0] == '/' {
		return fmt.Sprintf("%v%v", base, next) 
	}
	return fmt.Sprintf("%v/%v", base, next)
}
