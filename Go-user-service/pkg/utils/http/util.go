package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func ValidateRequestFormat(r *http.Request, expectedContentType string) error {
	contentType := r.Header.Get("Content-Type")
	if contentType != expectedContentType {
		return errors.New(fmt.Sprintf("Error: Content-Type header is not %s", expectedContentType))
	}
	return nil
}

func WriteResponse(data interface{}, w http.ResponseWriter, contentType string, status int) error {
	if resp, err := json.Marshal(data); err == nil {
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(status)
		_, err := w.Write(resp)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}
