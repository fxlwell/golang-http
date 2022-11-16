package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

var testPort = 8806

func TestMain(m *testing.M) {
	mux := http.NewServeMux()
	mux.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Query().Get("v"))
	})

	mux.HandleFunc("/post-form", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}

		fmt.Fprintf(w, r.FormValue("v"))
	})

	mux.HandleFunc("/post-body", func(w http.ResponseWriter, r *http.Request) {
		v, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Fprintf(w, "RequestBody err: %v", err)
			return
		}
		fmt.Fprintf(w, string(v))
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
