package msgo

import (
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

type routerGroup struct {
	name            string
	handleFuncMap   map[string]HandleFunc
	handleMethodMap map[string][]string
}

// func (g routerGroup) name()  {
//
// }
type router struct {
	routerGroups []*routerGroup
}

func (receiver *router) Group(name string) *routerGroup {
	group := routerGroup{
		name:          name,
		handleFuncMap: make(map[string]HandleFunc),
	}
	receiver.routerGroups = append(receiver.routerGroups, &group)
	return &group
}

func (r routerGroup) AddRouter(name string, handleFunc HandleFunc) {
	r.handleFuncMap[name] = handleFunc
}
func (r routerGroup) Any(name string, handleFunc HandleFunc) {
	r.handleFuncMap[name] = handleFunc
	r.handleMethodMap["ANY"] = append(r.handleMethodMap["ANY"], name)
}

// Engine  引擎
type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{
			routerGroups: []*routerGroup{},
		},
	}
}

func (e Engine) Run() {
	for _, g := range e.routerGroups {
		for key, val := range g.handleFuncMap {
			http.HandleFunc("/"+g.name+key, val)
		}
	}

	err := http.ListenAndServe(":8111", nil)

	if err != nil {
		log.Fatal(err)
	}
}
