package grains_api

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
)

func BindJSONStrict[T any](gc *gin.Context, out *T) error {
	decoder := json.NewDecoder(gc.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(out); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("multiple JSON objects are not allowed")
	}

	return nil
}
