package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	started := time.Now()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		data := time.Since(started).String()
		w.Write([]byte(data))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		duration := time.Since(started)
		if duration.Seconds() > 10 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("error:%v", duration.Seconds())))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		}
	})

	http.ListenAndServe(":8080", nil)

}
