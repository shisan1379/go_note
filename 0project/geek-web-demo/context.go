package geek_web_demo

import "net/http"

type Context struct {
	request  *http.Request
	response http.ResponseWriter
}
