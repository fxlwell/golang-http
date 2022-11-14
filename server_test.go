package http

import (
	"fmt"
	"net/http"
	"testing"
)

var testPort = 8806

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Query().Get("v"))
	})

	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, r.FormValue("v"))
	})

	mux.HandleFunc("/header", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.Header.Get("h"))
	})

	srv := DefaultServer
	srv.Addr = fmt.Sprintf("127.0.0.1:%d", testPort)
	srv.Handler = mux

	go srv.ListenAndServe()

	m.Run()
}
