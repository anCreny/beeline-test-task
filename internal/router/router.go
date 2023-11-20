package router

import (
	"log"
	"net/http"
	"strings"
)

type method string
type path string
type handleFunc func(http.ResponseWriter, *http.Request)

type Router struct {
	routsMap map[method]map[path]handleFunc
}

func New() *Router {
	return &Router{make(map[method]map[path]handleFunc)}
}

func (rr *Router) Handle(m method, p path, f handleFunc) {
	if _, found := rr.routsMap[m]; !found {
		rr.routsMap[m] = make(map[path]handleFunc)
	}

	rr.routsMap[m][p] = f
}

func (rr *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	m := method(r.Method)
	p := path(r.URL.Path)

	pathsMap, found := rr.routsMap[m]
	if !found {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var matchedPath path

	for path := range pathsMap {
		if matchPaths(string(path), string(p)) {
			matchedPath = path
		}
	}

	hFunc, found := pathsMap[matchedPath]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	hFunc(w, r)
}

func matchPaths(registeredPath string, requestPath string) bool {
	registeredPath = strings.Trim(registeredPath, "/")
	requestPath = strings.Trim(requestPath, "/")

	arrayRegPath := strings.Split(registeredPath, "/")
	arrayRequestPath := strings.Split(requestPath, "/")

	if len(arrayRegPath) != len(arrayRequestPath) {
		return false
	}

	for i := 0; i < len(arrayRegPath); i++ {
		if !isUrlParam(arrayRegPath[i]) && arrayRegPath[i] != arrayRequestPath[i] {
			return false
		}
	}

	return true
}

func isUrlParam(pathUnit string) bool {
	if len(pathUnit) < 3 {
		return false
	}

	if rune([]byte(pathUnit)[0]) == '{' && rune([]byte(pathUnit)[len(pathUnit)-1]) == '}' {
		return true
	}

	return false
}
