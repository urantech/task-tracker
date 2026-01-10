package task

import (
	"backend-api/internal/common"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(common.UserIdKey).(int64)
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

		msg, statusCode := common.HandleError(w, err)
		http.Error(w, msg, statusCode)

		return
	}

	common.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(common.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
		return
	}

	tasks, err := h.service.GetUserTasks(ctx, userId)
	if err != nil {
		msg, statusCode := common.HandleError(w, err)
		http.Error(w, msg, statusCode)
	}

	common.WriteJSON(w, http.StatusOK, tasks)
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId, ok := ctx.Value(common.UserIdKey).(int64)
	if !ok {
		http.Error(w, "Could not get user from context", http.StatusInternalServerError)
		return
	}

	taskIdParam := chi.URLParam(r, "id")

	taskId, err := validateTaskId(taskIdParam)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req UpdateRequest

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	resp, err := h.service.UpdateTask(ctx, taskId, userId, req)
	if err != nil {
		var ve *common.ValidationError
		if errors.As(err, &ve) {
			common.WriteJSON(w, http.StatusBadRequest, map[string]any{
				"errors": ve.Errors,
			})
			return
		}

		msg, statusCode := common.HandleError(w, err)
		http.Error(w, msg, statusCode)

		return
	}

	common.WriteJSON(w, http.StatusOK, resp)
}

func validateTaskId(taskIdParam string) (int64, error) {
	if taskIdParam == "" {
		return 0, fmt.Errorf("missing task ID")
	}

	id, err := strconv.ParseInt(taskIdParam, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid task ID: %w", err)
	}

	if id <= 0 {
		return 0, fmt.Errorf("task ID must be positive")
	}

	return id, nil
}
