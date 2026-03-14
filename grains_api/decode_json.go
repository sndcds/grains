package grains_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// DecodeJSONBody reads the request body, decodes JSON into 'out', and handles common errors
func DecodeJSONBody[T any](gc *gin.Context, apiRequest *Request) (out T, ok bool) {
	var zero T

	body, err := io.ReadAll(gc.Request.Body)
	if err != nil {
		apiRequest.Error(http.StatusBadRequest, "failed to read request body")
		return zero, false
	}
	if len(body) == 0 {
		apiRequest.Error(http.StatusBadRequest, "empty request body")
		return zero, false
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&out); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			apiRequest.Error(http.StatusBadRequest,
				fmt.Sprintf("invalid JSON syntax at offset %d", syntaxErr.Offset))
		case errors.As(err, &typeErr):
			field := typeErr.Field
			if field == "" {
				field = "(unknown)"
			}
			apiRequest.Error(http.StatusBadRequest,
				fmt.Sprintf("invalid type for field %q", field))
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			apiRequest.Error(http.StatusBadRequest, err.Error())
		default:
			apiRequest.Error(http.StatusBadRequest, err.Error())
		}
		return zero, false
	}

	if decoder.More() {
		apiRequest.Error(http.StatusBadRequest, "multiple JSON objects are not allowed")
		return zero, false
	}

	return out, true
}
