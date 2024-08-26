package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data any
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func (r JSON) Render(w http.ResponseWriter, code int) error {
	return WriteJSON(w, r.Data, code)
}
func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType[0])
}

func WriteJSON(w http.ResponseWriter, obj any, code int) error {
	writeContentType(w, jsonContentType[0])
	w.WriteHeader(code)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}
