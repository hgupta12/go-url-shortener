package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"
)

type APIServer struct{
	listenAddr string
	shortener Shortener
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func NewAPIServer(listenAddr string, shortener Shortener) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		shortener: shortener,
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)
		if !slices.Contains(methods, r.Method) {
			writeJSON(w, http.StatusMethodNotAllowed, ApiError{ Error: fmt.Sprintf("method not allowed - %s", r.Method)})
		} else if err := f(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}


func (s *APIServer) handleLongURL(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()
	request := new(Request)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		return err
	}

	hash := s.shortener.Shorten(request.Url)

	response := Response{Url: fmt.Sprintf("http://localhost%s/%s", s.listenAddr, hash)}
	
	return writeJSON(w, http.StatusOK, response)
}

func (s *APIServer) handleResolveShortURL(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	hash, ok := vars["hash"]
	if !ok {
		return fmt.Errorf("required parameter - hash missing")
	}

	url,err := s.shortener.Resolve(hash);
	if err != nil {
		log.Printf("failed to resolve hash %s: %s", hash, err)
		return writeJSON(w,http.StatusNotFound, ApiError{Error: err.Error()})
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	return nil
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	router.HandleFunc("/short", makeHTTPHandleFunc(s.handleLongURL, "POST"))
	router.HandleFunc("/{hash}", makeHTTPHandleFunc(s.handleResolveShortURL, "GET"))

	return http.ListenAndServe(s.listenAddr, router)
}