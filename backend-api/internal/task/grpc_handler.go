package task

import (
	taskv1 "backend-api/gen/go/task/v1"
	"context"
)

type GrpcHandler struct {
	taskv1.UnimplementedTaskAnalyticsServiceServer

	service *Service
}

func NewGrpcHandler(service *Service) *GrpcHandler {
	return &GrpcHandler{service: service}
}

func (h *GrpcHandler) CollectAndSendTaskAnalytics(
	ctx context.Context,
	req *taskv1.CollectAndSendTaskAnalyticsRequest,
) (*taskv1.CollectAndSendTaskAnalyticsResponse, error) {
	err := h.service.CollectAndSendTaskAnalytics(ctx)
	if err != nil {
		return &taskv1.CollectAndSendTaskAnalyticsResponse{
			Status: taskv1.JobStatus_JOB_STATUS_ERROR,
			Error:  err.Error(),
		}, nil
	}

	return &taskv1.CollectAndSendTaskAnalyticsResponse{
		Status: taskv1.JobStatus_JOB_STATUS_SUCCESS,
		Error:  "",
	}, nil
}
