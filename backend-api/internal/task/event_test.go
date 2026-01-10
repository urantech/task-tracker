package task

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageWriter struct {
	mock.Mock
}

func (m *MockMessageWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func TestProducer_ProduceDailyReport(t *testing.T) {
	tests := []struct {
		name        string
		report      DailyReportMsg
		mockError   error
		expectError bool
	}{
		{
			name: "successful send",
			report: DailyReportMsg{
				UserId:  1,
				Message: "You have 2 pending tasks.",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "kafka error",
			report: DailyReportMsg{
				UserId:  2,
				Message: "Error message",
			},
			mockError:   errors.New("kafka connection failed"),
			expectError: true,
		},
		{
			name: "correct message structure",
			report: DailyReportMsg{
				UserId:  3,
				Message: "Daily report for user 3",
			},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := &MockMessageWriter{}
			producer := &Producer{w: mockWriter}

			payload, _ := json.Marshal(tt.report)

			mockWriter.On("WriteMessages", mock.Anything, mock.MatchedBy(func(msgs []kafka.Message) bool {
				if len(msgs) != 1 {
					return false
				}
				msg := msgs[0]
				return string(msg.Key) == string([]byte{byte(tt.report.UserId)}) && string(msg.Value) == string(payload)
			})).Return(tt.mockError).Once()

			err := producer.ProduceDailyReport(context.Background(), tt.report)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "kafka publish error")
			} else {
				assert.NoError(t, err)
			}

			mockWriter.AssertExpectations(t)
		})
	}
}
