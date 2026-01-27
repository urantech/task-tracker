package job

import (
	"context"
	"log"

	taskv1 "api-proto/gen/go/task/v1"
)

type GrpcClient struct {
	client taskv1.TaskAnalyticsServiceClient
}

func NewGrpcClient(client taskv1.TaskAnalyticsServiceClient) *GrpcClient {
	return &GrpcClient{client: client}
}

func (c *GrpcClient) CollectAndSendTaskAnalytics(ctx context.Context) {
	if c.client == nil {
		log.Printf("gRPC client not initialized")
		return
	}

	req := &taskv1.CollectAndSendTaskAnalyticsRequest{}

	resp, err := c.client.CollectAndSendTaskAnalytics(ctx, req)
	if err != nil {
		log.Printf("Failed to collect and send task analytics: %v", err)
		return
	}

	log.Printf("Task analytics collected: status=%s, error=%s, total=%d, success=%d, failed=%d, failed_user_ids=%v",
		resp.Status, resp.Error, resp.TotalCount, resp.SuccessCount, resp.FailedCount, resp.FailedUserIds)
}
