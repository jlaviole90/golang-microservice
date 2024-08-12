package api 

import "net/http"

var (
	JSONSerializeFailureResponse   = []byte(`{"error": "Failed to serialize data to JSON"}`)
	JSONDeserializeFailureResponse = []byte(`{"error": "Failed to deserialize data from JSON"}`)
	IOErrorResponse                = []byte(`{"error": "Failed to read/write data/file"}`)
	InternalServerErrorResponse    = []byte(`{"error": "Internal server error"}`)
	BadRequestResponse             = []byte(`{"error": "Bad request"}`)
	ClientErrorResponse            = []byte(`{"error": "Client error"}`)
	NotFoundResponse               = []byte(`{"error": "Not found"}`)
)

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	// STATUS 500
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}

func ClientError(w http.ResponseWriter, error []byte, client string) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
	w.Header().Add("Client", client)
}

func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}

func NotFound(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusNotFound)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}

func IOError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}

func JSONSerializeError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}

func JSONDeserializeError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	_, e := w.Write(error)
	if e != nil {
		println("Failed to write response")
	}
}
