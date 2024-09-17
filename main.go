package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config := Config{}
	err := config.Read("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := RequestHandler{
		Config: config,
	}

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(handler.HandleRequest)
	http.ListenAndServe(config.GetListener(), r)
}
