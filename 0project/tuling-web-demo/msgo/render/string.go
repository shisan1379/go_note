package render

import (
	"fmt"
	"github.com/shisan1379/msgo/internal/bytesconv"
	"net/http"
)

type String struct {
	Format string
	Data   []any
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func (r String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, plainContentType[0])
}

func (r String) Render(w http.ResponseWriter) error {
	return WriteString(w, r.Format, r.Data)
}

func WriteString(w http.ResponseWriter, format string, data []any) (err error) {
	writeContentType(w, plainContentType[0])
	if len(data) > 0 {
		_, err = fmt.Fprintf(w, format, data...)
		return
	}
	_, err = w.Write(bytesconv.StringToBytes(format))
	return
}
