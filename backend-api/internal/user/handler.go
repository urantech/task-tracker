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

		msg, statusCode := common.HandleError(err)
		http.Error(w, msg, statusCode)

		return
	}

	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(common.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	user, err := h.service.GetUser(ctx, userId)
	if err != nil {
		msg, statusCode := common.HandleError(err)
		http.Error(w, msg, statusCode)

		return
	}

	common.WriteJSON(w, http.StatusOK, user)
}
