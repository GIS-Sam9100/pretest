package at

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

func GetSecretFromHeader(r *http.Request) (secret string) {
	if r.Header.Get("secret") != "" {
		secret = r.Header.Get("secret")
	} else if r.Header.Get("Secret") != "" {
		secret = r.Header.Get("Secret")
	}
	return
}

func GetLoginFromHeader(r *http.Request) (secret string) {
	if r.Header.Get("login") != "" {
		secret = r.Header.Get("login")
	} else if r.Header.Get("Login") != "" {
		secret = r.Header.Get("Login")
	}
	return
}

func Jsonstr(strc interface{}) string {
	jsonData, err := json.Marshal(strc)
	if err != nil {
		log.Fatal(err)
	}
	return string(jsonData)
}

func WriteJSON(respw http.ResponseWriter, statusCode int, content interface{}) {
	respw.Header().Set("Content-Type", "application/json")
	respw.WriteHeader(statusCode)
	respw.Write([]byte(Jsonstr(content)))
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// Set a max body size to prevent malicious attacks
	maxBytes := 1048576 // 1MB
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Read and decode the request body
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// Check if there is only one JSON value in the body
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single JSON value")
	}

	return nil
}

func WriteFile(w http.ResponseWriter, statusCode int, fileContent []byte) {
	w.Header().Set("Content-Disposition", "attachment; filename=\"file.ext\"")
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprint(len(fileContent)))
	w.WriteHeader(statusCode)
	w.Write(fileContent)
}

func WriteString(respw http.ResponseWriter, statusCode int, content string) {
	respw.WriteHeader(statusCode)
	respw.Write([]byte(content))
}
