package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	gofigure "github.com/common-nighthawk/go-figure"
	"github.com/gorilla/mux"
)

func main() {
	fmt.Print("\r\n")
	gofigure.NewColorFigure("LS AGENT RELAY", "epic", "green", true).Print()
	fmt.Print("\r\n\r\n\r\n\r\n")

	args := os.Args[1:]
	if len(args) != 1 {
		envv := os.Getenv("LSAGENTRELAY_CONFIG")
		if envv != "" {
			args = []string{envv}
		} else {
			fmt.Println("Usage: lsagentrelay [config.yaml path]")
			fmt.Println("Using default config.yaml")
			args = []string{"config.yaml"}
		}
	}

	config := Config{}
	err := config.Read(args[0])
	if err != nil {
		fmt.Println("Configuration file could not be read: ", err)
		return
	}

	debugenv := strings.ToLower(os.Getenv("LSAGENTRELAY_DEBUG"))
	if debugenv != "" {
		config.Listen.Debug = debugenv == "true" || debugenv == "1" || debugenv == "yes"
	}

	handler := RequestHandler{
		Config:   config,
		DebugLog: func(message string, args ...any) {},
	}

	if config.Listen.Debug {
		fmt.Println("!! Debug logging enabled !!")
		handler.DebugLog = func(message string, args ...any) {
			fmt.Printf(message, args...)
		}
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
