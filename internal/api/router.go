package api

import (
	"encoding/json"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", fooHandler)
	return mux
}

func fooHandler(w http.ResponseWriter, r *http.Request) {
	res, _ := json.Marshal(
		map[string]string{
			"body": "Hey from Server",
		},
	)
	w.Write(res)
}
