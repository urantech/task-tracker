package task

import (
	taskv1 "api-proto/gen/go/task/v1"
	"context"
	"fmt"
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
	stats, err := h.service.CollectAndSendTaskAnalytics(ctx)
	if err != nil {
		return &taskv1.CollectAndSendTaskAnalyticsResponse{
			Status: taskv1.JobStatus_JOB_STATUS_ERROR,
			Error:  err.Error(),
		}, nil
	}

	status := taskv1.JobStatus_JOB_STATUS_SUCCESS

	var errMsg string

	if len(stats.FailedUserIds) > 0 {
		errMsg = fmt.Sprintf("failed to send %d reports", len(stats.FailedUserIds))

		if stats.SuccessCount == 0 {
			status = taskv1.JobStatus_JOB_STATUS_ERROR
		}
	}

	return &taskv1.CollectAndSendTaskAnalyticsResponse{
		Status:        status,
		Error:         errMsg,
		TotalCount:    int32(stats.Total),
		SuccessCount:  int32(stats.SuccessCount),
		FailedCount:   int32(len(stats.FailedUserIds)),
		FailedUserIds: stats.FailedUserIds,
	}, nil
}
