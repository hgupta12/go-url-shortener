package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"fmt"
	"log"
	"net/http"
)

type APIServer struct{
	listenAddr string
	shortener Shortener
}
func NewAPIServer(listenAddr string, shortener Shortener) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		shortener: shortener,
	}
}

func (s *APIServer) handleLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "method not supported - %s", r.Method)
		return
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	request := new(Request)
	err := decoder.Decode(request)

	if err != nil {
		fmt.Fprintln(w, "Error decoding JSON")
		log.Fatal(err)
	}

	hash := s.shortener.Shorten(request.Url)

	response := Response{Url: fmt.Sprintf("http://localhost%s/%s", s.listenAddr, hash)}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(&response)

	if err!= nil {
		fmt.Fprintln(w, "Error")
		log.Fatal(err)
	}	
}

func (s *APIServer) handleResolveShortURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash, ok := vars["hash"]
	if !ok {
		fmt.Fprintln(w, "missing required paramter - hash")
	}

	url, err := s.shortener.Resolve(hash)
	if err != nil {
		fmt.Fprintln(w, err)
		log.Fatal(err)
	}

	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/short", s.handleLongURL)
	router.HandleFunc("/{hash}", s.handleResolveShortURL)

	return http.ListenAndServe(s.listenAddr, router)
}