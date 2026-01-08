package auth

import (
	"backend-api/internal/common"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(r.Context(), req)
	if err != nil {
		var ve *common.ValidationError
		if errors.As(err, &ve) {
			common.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errors": ve.Errors,
			})
			return
		}

		common.HandleError(w, err)

		return
	}

	var resp LoginResponse

	resp.AccessToken = token

	common.WriteJSON(w, http.StatusOK, resp)
}
