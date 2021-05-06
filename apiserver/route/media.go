package api

import "net/http"

func MediaHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all dogs"))
}
