package geek_web_demo

type router struct {
	Trees map[string]*node
}

type node struct {
	Path     string
	Children map[string]*node
	Handler  HandleFunc
}

func NewRouter() *router {
	return &router{Trees: make(map[string]*node)}
}
func (h *router) AddRoute(method string, pattern string, handler HandleFunc) {
	//http.HandleFunc(pattern, handler)
}
