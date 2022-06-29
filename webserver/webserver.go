package webserver

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
)

var kvs = make(map[string]string)

func Run() {
	http.HandleFunc("/store", storeHandler)
	http.ListenAndServe(":8081", nil)
}

func storeHandler(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		writeErrorResponse(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	switch r.Method {
		case http.MethodGet:
			var req readRequest
			err := decoder.Decode(&req)
			if err != nil {
				msg := getJSONErrorMesssage(err)
				writeErrorResponse(w, msg, http.StatusBadRequest)
				return
			}

			val, ok := read(&kvs, req.Key)
			if !ok {
				writeSuccessResponse(w, "Key not found", http.StatusOK)
			} else {
				writeSuccessResponse(w, val, http.StatusOK)
			}

		case http.MethodPost:
			var req writeRequest
			err := decoder.Decode(&req)
			if err != nil {
				msg := getJSONErrorMesssage(err)
				writeErrorResponse(w, msg, http.StatusBadRequest)
				return
			}

			write(&kvs, req.Key, req.Val)
			writeSuccessResponse(w, req.Val, http.StatusOK)
		case http.MethodDelete:
			var req deleteRequest
			err := decoder.Decode(&req)
			if err != nil {
				msg := getJSONErrorMesssage(err)
				writeErrorResponse(w, msg, http.StatusBadRequest)
				return
			}

			del(&kvs, req.Key)
			writeSuccessResponse(w, "", http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getJSONErrorMesssage(err error) string {
	var unmarshalErr *json.UnmarshalTypeError
	if errors.As(err, &unmarshalErr) {
		return "Bad Request. Wrong Type provided for field "+unmarshalErr.Field
	}
	return "Bad Request, JSON error: "+err.Error()
}

func writeErrorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["op_succeeded"] = "false"
	resp["error"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func writeSuccessResponse(w http.ResponseWriter, data string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["op_succeeded"] = "true"
	resp["response"] = data
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

type readRequest struct {
	Key string `json:key`
}

type writeRequest struct {
	Key string `json:key`
	Val string `json:val`
}

type deleteRequest struct {
	Key string `json:key`
}

func del(kvs *map[string]string, k string) {
	delete(*kvs, k)
}

func write(kvs *map[string]string, k string, v string) {
	(*kvs)[k] = v
}

func read(kvs *map[string]string, k string) (string, bool) {
	v, ok := (*kvs)[k]
	return v, ok
}

func printMap(m *map[string]string) {
	for k, v := range *m {
        fmt.Println(k, ":", v)
    }
}


//  curl -X POST localhost:8081/store -H 'Content-type: application/json' -d '{"key":"test_key", "val":"test_val"}'
//  curl -X GET localhost:8081/store -H 'Content-type: application/json' -d '{"key":"test_key"}'
//  curl -X DELETE localhost:8081/store -H 'Content-type: application/json' -d '{"key":"test_key"}'
