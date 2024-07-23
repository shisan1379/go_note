package msgo

import (
	"log"
	"net/http"
)

type HandleFunc func(w http.ResponseWriter, r *http.Request)

type router struct {
	handleMap map[string]HandleFunc
}

func (r router) Add(name string, handleFunc HandleFunc) {
	r.handleMap[name] = handleFunc
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{
			handleMap: make(map[string]HandleFunc),
		},
	}
}

func (e Engine) Run() {
	err := http.ListenAndServe(":8111", nil)

	if err != nil {
		log.Fatal(err)
	}
}
