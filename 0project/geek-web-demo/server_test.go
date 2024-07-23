package geek_web_demo

import (
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	h := &HttpServer{}
	//方式一
	http.ListenAndServe(":8080", h)

	//h.Get()

	//方式二
	h.Start(":8081")
}
