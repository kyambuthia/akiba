package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const maxRequestBodyBytes int64 = 1 << 20 // 1MB

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		switch {
		case errors.As(err, &maxErr):
			writeError(w, http.StatusBadRequest, "bad_request", "request body too large", nil)
		default:
			writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload", nil)
		}
		return false
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload", nil)
		return false
	}
	return true
}
