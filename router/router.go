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

	apicollectionshelp := apicollections
	rt.setroute(apicollectionshelp, "GET", func(w http.ResponseWriter,  r *http.Request) {
		w.Write([]byte("#documentation for interacting with collections"))
	})
	apicollectionsone := appendPath(apicollections, "*")
	rt.setroute(apicollectionsone, "GET", controllers.GetCollectionHandler)
	rt.setroute(apicollectionsone, "POST", controllers.PostCollectionHandler)
	rt.setroute(apicollectionsone, "PUT", controllers.PutInCollectionHandler)
	rt.setroute(apicollectionsone, "DELETE", controllers.DeleteCollectionHandler)
	apicollectionsurlsone := appendPath(apicollectionsone, "urls/*")
	rt.setroute(apicollectionsurlsone, "DELETE", controllers.DeleteURLFromCollection)
	return rt
}

type router map[int][]routeEntry

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedHandler := http.NotFound
	reqPathTokens := controllers.ParsePathTokens(r.URL.Path)
	bestScore := 0
	for _, rentry := range (*rt)[len(reqPathTokens)] {
		if rentry.doesMatchMethod(r) {
			score := rentry.doesMatchPath(reqPathTokens) 
			if score > bestScore { 
				bestScore = score
				matchedHandler = rentry.Handler.ServeHTTP
			}
		}
	}

	matchedHandler(w, r)
}

func (rentry *routeEntry) doesMatchPath(requestPathTokens []string) (int) {
	scr := 0 
	fmt.Println(rentry.Path)
	fmt.Println(requestPathTokens)
	for i, val := range requestPathTokens {
		currToken := rentry.Path[i]
		if val != currToken && currToken != "*" {
			return -1
		}
		if val == currToken {
			scr++
		}
	}
	return scr 
}

func (rt *router) setroute(p string, m string, hf http.HandlerFunc) {
	pathTokens := controllers.ParsePathTokens(p) 
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

func appendPath(base, next string) string {
	if next[0] == '/' {
		return fmt.Sprintf("%v%v", base, next) 
	}
	return fmt.Sprintf("%v/%v", base, next)
}
