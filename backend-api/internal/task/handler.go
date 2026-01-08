package task

import (
	"backend-api/internal/auth"
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

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(auth.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
		return
	}

	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.CreateTask(ctx, req, userId)
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
