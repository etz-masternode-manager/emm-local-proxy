package emm_local_proxy

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

type ResponseError struct {
	Error string `json:"error"`
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

		w.WriteHeader(response.StatusCode)

		io.Copy(w, response.Body)
	}

}

func main() {
	InitializePublicKey()

	http.HandleFunc("/", processRequest)
	log.Fatal(http.ListenAndServe(":34521", nil))
}
