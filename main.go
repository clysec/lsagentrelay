package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("Usage: lsagentrelay [config.yaml path]")
		fmt.Println("Using default config.yaml")
		args = []string{"config.yaml"}
	}

	config := Config{}
	err := config.Read(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	handler := RequestHandler{
		Config: config,
	}

	r := mux.NewRouter()
	r.Path("/lsagent").HandlerFunc(handler.HandleRequest)

	if config.Listen.Tls.Enabled {
		cert, err := tls.LoadX509KeyPair(config.Listen.Tls.Cert, config.Listen.Tls.Key)
		if err != nil {
			fmt.Println(err)
			return
		}

		server := &http.Server{
			Addr:    config.GetListener(),
			Handler: r,
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{cert},
			},
		}

		go func() {
			fmt.Println("Starting TLS server on", config.GetListener())
			err := server.ListenAndServeTLS("", "")
			if err != nil {
				fmt.Println(err)
			}
		}()
	} else {
		fmt.Println("Starting HTTP server on", config.GetListener())
		http.ListenAndServe(config.GetListener(), r)
	}
}
