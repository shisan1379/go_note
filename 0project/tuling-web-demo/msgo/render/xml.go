package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data any
}

var xmlContentType = []string{"application/xml; charset=utf-8"}

func (r XML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)
	return xml.NewEncoder(w).Encode(r.Data)
}

func (r XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType[0])
}
