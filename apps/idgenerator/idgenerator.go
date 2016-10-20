package main

import (
	"fmt"
	"net/http"

	"github.com/rs/xid"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in f", r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		id := xid.New()
		fmt.Fprintf(w, "%s", id.String())
	})

	http.ListenAndServe(":8888", nil)
}
