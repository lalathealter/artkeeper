package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/lalathealter/artkeeper/controllers"
)

const (
	staticfilesdir = "static"
	staticfiles = "/static/*+"
	apiurls        = "/api/urls"
	apicollections = "/api/collections"
	apiusers = "/api/users"
	apisession = "/api/session"
)

func Use() *router {
	rt := &router{	}

	rt.setroute("/", "GET", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world!"))
	})

	rt.setroute(staticfiles, "GET", func(w http.ResponseWriter, r *http.Request) {
		statFS := http.StripPrefix("/" + staticfilesdir, http.FileServer(http.Dir(staticfilesdir)))
		statFS.ServeHTTP(w, r)
	})
	
	apiurlshelp := apiurls 
	rt.setroute(apiurlshelp, "GET", controllers.HelpURLHandler)
	apiurlsone := appendPath(apiurls, "*") 
	rt.setroute(apiurlsone, "POST", controllers.PostURLhandler)
	rt.setroute(apiurlsone, "GET", controllers.GetURLHandler)
	rt.setroute(apiurlsone, "DELETE", controllers.DeleteURLHandler)

	apiurlslatest := appendPath(apiurls, "latest") 
	rt.setroute(apiurlslatest, "GET", controllers.GetLatestURLsHandler)

	rt.setroute(apicollections, "POST", controllers.PostCollectionHandler)
	apicollectionshelp := apicollections
	rt.setroute(apicollectionshelp, "GET", controllers.HelpCollectionHandler)

	apicollectionsone := appendPath(apicollections, "*")
	rt.setroute(apicollectionsone, "GET", controllers.GetCollectionHandler)
	rt.setroute(apicollectionsone, "DELETE", controllers.DeleteCollectionHandler)

	apicollectionsurls := appendPath(apicollectionsone, "urls") 
	rt.setroute(apicollectionsurls, "PUT", controllers.PutInCollectionHandler)
	rt.setroute(apicollectionsurls, "GET", controllers.GetURLsFromCollectionHandler)
	apicollectionsurlsone := appendPath(apicollectionsone, "urls/*")
	rt.setroute(apicollectionsurlsone, "DELETE", controllers.DeleteURLFromCollection)

	apicollectionstags := appendPath(apicollectionsone, "tags/*")
	rt.setroute(apicollectionstags, "PUT", controllers.AttachTagToCollectionHandler)
	rt.setroute(apicollectionstags, "DELETE", controllers.DetachTagFromCollectionHandler)

	apiusersnew := appendPath(apiusers, "new")
	rt.setroute(apiusersnew, "POST", controllers.UserRegistrationHandler)
	rt.setroute(apiusers, "POST", controllers.UpdateUserHandler)

	rt.setroute(apisession, "POST", controllers.PostSessionHandler)
	return rt
}

type router map[int]routeEntriesMap

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	matchedHandler := rt.matchHandlerFor(r)
	matchedHandler(w, r)
}

func (rt *router) matchHandlerFor(r *http.Request) http.HandlerFunc {
	reqPathTokens := controllers.ParsePathTokens(r.URL.Path)
	reqLen := len(reqPathTokens)
	method := r.Method
	exactRentriesMap := (*rt)[reqLen] 
	var matchedHandler http.HandlerFunc 

	exactRoutes, isMethodAllowed := exactRentriesMap[method]
	matchedHandler = findBestHandler(reqPathTokens, exactRoutes)
	if matchedHandler != nil {
		return matchedHandler
	}

	dynamicRoutes := (*rt)[0][method]
	matchedHandler = findBestHandler(reqPathTokens, dynamicRoutes)
	if matchedHandler != nil {
		return matchedHandler
	}

	if !isMethodAllowed {
		return exactRentriesMap.methodNotAllowedHandler
	}

	return http.NotFound
}

func findBestHandler(reqPathTokens []string, routes []*routeEntry) http.HandlerFunc {
	var matchedHandler http.HandlerFunc
	bestScore := 0

	for _, rentry := range routes {
		score := rentry.doesMatchPath(reqPathTokens)  
		if score > bestScore { 
			bestScore = score
			matchedHandler = rentry.Handler.ServeHTTP
		}
	}

	return matchedHandler
}

type routeEntriesMap map[string][]*routeEntry 

func (rens routeEntriesMap) methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	allowedMethods := rens.getAllowedMethods()  
	allowedMethodsHeaderString := strings.Join(allowedMethods, ", ")
	w.Header().Set("Allow", allowedMethodsHeaderString)
	w.WriteHeader(http.StatusMethodNotAllowed)
}


func (rens routeEntriesMap) getAllowedMethods() []string {
	methods := make([]string, len(rens))
	i := 0
	for m := range rens {
		methods[i] = m
		i++
	}
	return methods 
}

func (rentry *routeEntry) doesMatchPath(requestPathTokens []string) (int) {
	scr := 0 
	for i:= 0; i < len(rentry.Path) && i < len(requestPathTokens); i++ {
		reqToken := requestPathTokens[i]
		routeToken := rentry.Path[i]

		var isParameter bool 
		if (routeToken != "") {
			isParameter = routeToken[0] == '*'
		} 

		if reqToken != routeToken && !isParameter {
			return -1
		}
		if reqToken == routeToken {
			scr++
		}
	}
	return scr 
}

func (rt *router) setroute(p string, m string, hf http.HandlerFunc) {
	pathTokens := controllers.ParsePathTokens(p) 
	rentry := &routeEntry{
		Path:    pathTokens,
		Method:  strings.ToUpper(m),
		Handler: hf,
	}
	l := len(pathTokens)
	rt.appendRentry(l, rentry)

	lastToken := pathTokens[l - 1]
	if lastToken != "" {
		if lastToken[len(lastToken) - 1] == '+' {
			rt.appendRentry(0, rentry)
		}
	}
}

func (rt *router) appendRentry(size int, rentry *routeEntry) {
	method := rentry.Method
	if _, ok := (*rt)[size]; !ok {
		(*rt)[size] = make(routeEntriesMap)
	} 
	(*rt)[size][method] = append((*rt)[size][method], rentry)
}

type routeEntry struct {
	Path    []string
	Method  string
	Handler http.HandlerFunc
}



func appendPath(base string, nextOnes ...string) string {
	result := base
	for _, next := range nextOnes {
		if next[0] == '/' {
			result = fmt.Sprintf("%v%v", result, next) 
		} else {
			result = fmt.Sprintf("%v/%v", result, next)
		}
	}
	return result
}
