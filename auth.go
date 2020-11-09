package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("host:", r.Host)
		log.Println("method:", r.Method)
		log.Println("path:", r.URL.RawPath)
		log.Println("headers:", r.Header)
		body, err := ioutil.ReadAll(r.Body)
		if err == nil {
			log.Println("body:", string(body))
			r.Body.Close()
		} else if err == io.EOF {
			log.Println("NO REQUEST BODY")
		}
		w.Header().Set("Authorization", "Bearer BIGJWT")
		w.WriteHeader(200)
	}))
}
