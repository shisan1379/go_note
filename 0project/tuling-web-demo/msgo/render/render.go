package render

import "net/http"

// 所有 http 返回数据 都实现此接口
type Render interface {
	// 发送数据数据
	Render(w http.ResponseWriter) error
	// 设置响应头 - 也就是数据类型
	WriteContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, val string) {
	w.Header().Set("Content-Type", val)
}
