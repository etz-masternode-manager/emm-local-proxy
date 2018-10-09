package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type GethRequest struct {
	Port    int    `json:"port"`
	Request string `json:"request"`
}

type Response struct {
	Result string `json:"result"`
}

type ResponseError struct {
	Error string `json:"error"`
}

func responseResult(w io.Writer, result []byte) error {
	/*encrypted, err := encrypt(result)
	if err != nil {
		return responseError(w, err.Error())
	}
	*/
	var response Response
	response.Result = string(result) // base64.StdEncoding.EncodeToString(encrypted)
	return json.NewEncoder(w).Encode(response)
}

func responseError(w io.Writer, err string) error {
	var response ResponseError
	response.Error = err
	return json.NewEncoder(w).Encode(response)
}

func processRequest(w http.ResponseWriter, r *http.Request) {
	encodedSignature := r.Header.Get("Signature")
	if len(encodedSignature) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		responseError(w, "Empty signature")
		return
	}

	signature, err := base64.StdEncoding.DecodeString(encodedSignature)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError(w, "Can't parse signature: "+err.Error())
		return
	}

	requestType := r.Header.Get("Request-Type")
	if len(requestType) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		responseError(w, "Empty request type")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		responseError(w, "Can't read body: "+err.Error())
		return
	}

	if !VerifySignature(body, signature) {
		w.WriteHeader(http.StatusBadRequest)
		responseError(w, "Wrong signature")
		return
	}

	switch requestType {
	case "geth":
		var gethRequest GethRequest
		err := json.NewDecoder(bytes.NewReader(body)).Decode(&gethRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responseError(w, "Can't parse request: "+err.Error())
			return
		}

		response, err := http.Post("http://localhost:"+strconv.Itoa(gethRequest.Port), "application/json", strings.NewReader(gethRequest.Request))
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			responseError(w, "Can't perform Geth request: "+err.Error())
			return
		}

		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			responseError(w, "Request performed, buy can't read Geth response: "+err.Error())
			return
		}

		responseResult(w, responseBody)
	case "masternode":
		var mnRequest MasternodeRequest
		err := json.NewDecoder(bytes.NewReader(body)).Decode(&mnRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			responseError(w, "Can't parse request: "+err.Error())
			return
		}

		response, err := MasternodeRequestProcess(mnRequest)

		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			responseError(w, "Can't perform Geth request: "+err.Error())
			return
		}

		responseResult(w, []byte(response))
	}
}

func main() {
	InitializePublicKey()

	http.HandleFunc("/", processRequest)
	log.Fatal(http.ListenAndServe(":34521", nil))
}
