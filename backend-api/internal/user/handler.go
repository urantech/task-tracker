package user

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

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.RegisterUser(r.Context(), req)
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

	common.WriteJSON(w, http.StatusCreated, resp)
}
