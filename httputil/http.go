package httputil

import (
	"encoding/json"
	"io"
	"net/http"
)

func ParseBody[T any](body io.ReadCloser, target *T, responseWriter http.ResponseWriter) error {
	defer body.Close()
	bytes, errRead := io.ReadAll(body)
	if errRead != nil {
		WritePlainError(responseWriter, http.StatusInternalServerError)
		return errRead
	}

	if errUnmarshal := json.Unmarshal(bytes, target); errUnmarshal != nil {
		WritePlainError(responseWriter, http.StatusInternalServerError)
		return errUnmarshal
	}

	return nil
}

func WritePlainError(responseWriter http.ResponseWriter, statusCode int) {
	http.Error(responseWriter, http.StatusText(statusCode), statusCode)
}
