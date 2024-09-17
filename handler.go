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
	Config Config
}

func (h *RequestHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
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
		http.Error(w, "Error parsing body", http.StatusBadRequest)
		return
	}

	if len(parsedBody.Value["Action"]) == 0 {
		http.Error(w, "No action found", http.StatusBadRequest)
		return
	}

	action := parsedBody.Value["Action"][0]

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
		http.Error(w, "Action not found", http.StatusBadRequest)
	}
}

func (h *RequestHandler) HandleHello(w http.ResponseWriter, r *http.Request) {
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

	w.Write(body)
}
