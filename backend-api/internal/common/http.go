package common

import (
	"encoding/json"
	"errors"
	"net/http"
)

const internalServerError = "internal server error"

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "encode to json error", http.StatusInternalServerError)
		return
	}
}

func HandleError(err error) (string, int) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return err.Error(), http.StatusBadRequest

	case errors.Is(err, ErrUserAlreadyExists):
		return err.Error(), http.StatusConflict

	case errors.Is(err, ErrInvalidCredentials):
		return err.Error(), http.StatusUnauthorized

	case errors.Is(err, ErrTaskNotFound):
		return err.Error(), http.StatusNotFound

	default:
		return internalServerError, http.StatusInternalServerError
	}
}
