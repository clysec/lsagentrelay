package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/clysec/greq"
)

type RequestHandler struct {
	Config   Config
	DebugLog func(message string, args ...any)
}

func (h *RequestHandler) RequestToDisplayString(r *http.Request, parsedBody *multipart.Form, rawBody []byte) string {
	headerString := ""
	for k, v := range r.Header {
		headerString += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	bodyString := ""
	if len(parsedBody.Value) > 0 {
		for k, v := range parsedBody.Value {
			bodyString += fmt.Sprintf("%s: %s\r\n", k, v)
		}
	} else {
		bodyString = string(rawBody)
	}

	return fmt.Sprintf("%s %s [src %s]\r\n%s\r\n\r\n%s", r.Method, r.URL.Path, r.RemoteAddr, headerString, bodyString)
}

func (h *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil || len(bodyBytes) == 0 {
		h.DebugLog("Error reading body from client %s: %v", r.RemoteAddr, err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var boundary string

	if strings.Contains(contentType, "boundary") {
		boundary = strings.Split(contentType, "boundary=")[1]
	} else {
		if bytes.Contains(bodyBytes, []byte("boundary=")) {
			boundary = strings.Split(string(bodyBytes), "boundary=")[1]
			boundary = strings.Split(boundary, "\r\n")[0]
		} else {
			firstLine := strings.Split(string(bodyBytes), "\r\n")[0]
			if strings.Contains(firstLine, "-----") {
				boundary = firstLine
			} else {
				h.DebugLog("Error reading body from client %s: error reading boundary, first line is %s", r.RemoteAddr, firstLine)
				http.Error(w, "No boundary found", http.StatusBadRequest)
				return
			}
		}
	}

	if len(boundary) == 44 {
		boundary = strings.TrimPrefix(boundary, "--")
	}

	parsedBody, err := multipart.NewReader(bytes.NewReader(bodyBytes), boundary).ReadForm(10 << 20)
	if err != nil {
		h.DebugLog("Error parsing body from client %s: %s", r.RemoteAddr, err)
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	if len(parsedBody.Value["Action"]) == 0 {
		h.DebugLog("No action found, body: %s", string(bodyBytes))
		http.Error(w, "No action found", http.StatusBadRequest)
		return
	}

	action := parsedBody.Value["Action"][0]
	h.DebugLog("Received %s request from %s: %s", action, r.RemoteAddr, h.RequestToDisplayString(r, parsedBody, bodyBytes))

	switch action {
	case "Hello":
		h.HandleHello(w, r)
	case "AssetStatus":
		h.HandleAssetStatus(w, r, parsedBody, bodyBytes)
	case "Config":
		h.HandleConfig(w, r, parsedBody, bodyBytes)
	case "ScanData":
		h.HandleScanData(w, r, parsedBody, bodyBytes)
	default:
		h.DebugLog("Received unknown action %s from %s", action, r.RemoteAddr)
		http.Error(w, "Action not found", http.StatusBadRequest)
	}
}

func (h *RequestHandler) HandleHello(w http.ResponseWriter, r *http.Request) {
	h.DebugLog("Sending OK response to %s", r.RemoteAddr)
	fmt.Fprintf(w, "OK")
}

func (h *RequestHandler) HandleAssetStatus(w http.ResponseWriter, r *http.Request, bodyData *multipart.Form, rawBody []byte) {
	req := greq.PostRequest(h.Config.GetLansweeperUrl()).
		WithHeaders(map[string]interface{}{
			"Content-Type": r.Header["Content-Type"][0],
		}).
		WithByteBody(rawBody)

	if h.Config.Lansweeper.IgnoreSSL {
		req = req.TlsSetNovalidate()
	}

	resp, err := req.Execute()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error calling agent", http.StatusInternalServerError)
		return
	}

	body, err := resp.BodyBytes()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	h.DebugLog("Sending response to %s", r.RemoteAddr)
	w.Write(body)
}

func (h *RequestHandler) HandleConfig(w http.ResponseWriter, r *http.Request, bodyData *multipart.Form, rawBody []byte) {
	req := greq.PostRequest(h.Config.GetLansweeperUrl()).
		WithHeaders(map[string]interface{}{
			"Content-Type": r.Header["Content-Type"][0],
		}).
		WithByteBody(rawBody)

	if h.Config.Lansweeper.IgnoreSSL {
		req = req.TlsSetNovalidate()
	}

	resp, err := req.Execute()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error calling agent", http.StatusInternalServerError)
		return
	}

	body, err := resp.BodyBytes()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	replaced := h.Config.CompiledRegexp.ReplaceAll(body, []byte(h.Config.Rewrite.ProxyHostname))
	h.DebugLog("Sending response to %s", r.RemoteAddr)
	w.Write(replaced)
}
func (h *RequestHandler) HandleScanData(w http.ResponseWriter, r *http.Request, bodyData *multipart.Form, rawBody []byte) {
	req := greq.PostRequest(h.Config.GetLansweeperUrl()).
		WithHeaders(map[string]interface{}{
			"Content-Type": r.Header["Content-Type"][0],
		}).
		WithByteBody(rawBody)

	if h.Config.Lansweeper.IgnoreSSL {
		req = req.TlsSetNovalidate()
	}

	resp, err := req.Execute()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error calling agent", http.StatusInternalServerError)
		return
	}

	body, err := resp.BodyBytes()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Error reading body", http.StatusInternalServerError)
		return
	}

	h.DebugLog("Sending response to %s", r.RemoteAddr)
	w.Write(body)
}
