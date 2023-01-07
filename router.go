/*****************************************************************************
 * router.go
 * Name: Arsh Banerjee
 * NetId: arshb
 *****************************************************************************/

package http_router

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Student defined types or constants go here

type path struct {
	Handler       http.HandlerFunc
	HTTPMethod    string
	Name          *regexp.Regexp
	CapturePhrase []string
	OriginalPath  string
}

// HTTPRouter stores the information necessary to route HTTP requests
type HTTPRouter struct {
	Routes []path
}

// NewRouter creates a new HTTP Router, with no initial routes
func NewRouter() *HTTPRouter {
	return new(HTTPRouter)
}

// AddRoute adds a new route to the router, associating a given method and path
// pattern with the designated http handler.
func (router *HTTPRouter) AddRoute(method string, pattern string, handler http.HandlerFunc) {
	OP := pattern
	var Phrase []string
	Phrase = append(Phrase, "-1")
	capture := regexp.MustCompile(":([^/]+)/?")

	if len(capture.FindStringSubmatch(pattern)) > 0 {
		for i, array := range capture.FindAllStringSubmatch(pattern, -1) {
			if i == 0 {
				Phrase[i] = array[1]
			} else {
				Phrase = append(Phrase, array[1])
			}
		}
		pattern = capture.ReplaceAllString(pattern, "([^/]+)/?")
	}

	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	for i, route := range router.Routes {
		if route.Name.String() == ("^"+pattern+"$") && route.HTTPMethod == method {
			fmt.Println("Running")
			if route.CapturePhrase[0] != "-1" && Phrase[0] == "-1" {
				router.Routes = append(router.Routes, path{Handler: handler, HTTPMethod: method, Name: regexp.MustCompile("^" + pattern + "$"), CapturePhrase: Phrase, OriginalPath: OP})
				return
			}

			router.Routes[i].Handler = handler
			router.Routes[i].CapturePhrase = Phrase
			router.Routes[i].OriginalPath = OP
			return
		}
	}

	router.Routes = append(router.Routes, path{Handler: handler, HTTPMethod: method, Name: regexp.MustCompile("^" + pattern + "$"), CapturePhrase: Phrase, OriginalPath: OP})
	return
}

// ServeHTTP writes an HTTP response to the provided response writer
// by invoking the handler associated with the route that is appropriate
// for the provided request.
func (router *HTTPRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	pattern := strings.TrimPrefix(request.URL.Path, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	var bestRoute path
	found := false
	print := ""
	bestfound := false

	for _, route := range router.Routes {

		if route.Name.MatchString(pattern) {

			if route.CapturePhrase[0] != "-1" {
				if found {
					i := 0
					bestfound = false
					string1 := route.OriginalPath
					string2 := bestRoute.OriginalPath
					fmt.Println(string1)
					fmt.Println(string2)
					for !bestfound {

						if i == len(bestRoute.CapturePhrase) {
							break
						}

						if i == len(route.CapturePhrase) {
							bestRoute = route
							break
						}

						routeCapIndex := strings.Index(string1, route.CapturePhrase[i])
						bestrouteCapIndex := strings.Index(string2, bestRoute.CapturePhrase[i])
						routeCount := strings.Count(string1[:routeCapIndex], "/")
						bestrouteCount := strings.Count(string2[:bestrouteCapIndex], "/")
						if routeCount > bestrouteCount {
							bestRoute = route
							bestfound = true
						}
						if routeCount < bestrouteCount {
							bestfound = true
						}
						string1 = string1[routeCapIndex+len(route.CapturePhrase[i]):]
						string2 = string2[bestrouteCapIndex+len(bestRoute.CapturePhrase[i]):]
						i = i + 1
					}
				} else {
					bestRoute = route
					found = true
				}

			}

			if route.HTTPMethod == request.Method && route.CapturePhrase[0] == "-1" {
				route.Handler(response, request)
				return
			}
		}
	}

	for i, name := range bestRoute.CapturePhrase {
		if i != 0 {
			print = print + "&"
		}
		print = print + name + "=" + bestRoute.Name.FindStringSubmatch(pattern)[i+1]
	}

	if found && bestRoute.HTTPMethod == request.Method {
		fmt.Println(print)
		request.URL.RawQuery = print
		bestRoute.Handler(response, request)
		return
	}

	http.NotFound(response, request)
	return
}
