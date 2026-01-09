package user

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

func TestProducer_ProduceRegistration(t *testing.T) {
	tests := []struct {
		name        string
		user        UserResponse
		mockError   error
		expectError bool
	}{
		{
			name: "successful send",
			user: UserResponse{
				Id:    1,
				Email: "test@example.com",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name: "kafka error",
			user: UserResponse{
				Id:    2,
				Email: "error@example.com",
			},
			mockError:   errors.New("kafka connection failed"),
			expectError: true,
		},
		{
			name: "correct message structure",
			user: UserResponse{
				Id:    3,
				Email: "struct@example.com",
			},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockWriter := &MockMessageWriter{}
			producer := &Producer{w: mockWriter}

			payload, _ := json.Marshal(tt.user)

			mockWriter.On("WriteMessages", mock.Anything, mock.MatchedBy(func(msgs []kafka.Message) bool {
				if len(msgs) != 1 {
					return false
				}
				msg := msgs[0]
				return string(msg.Key) == tt.user.Email && string(msg.Value) == string(payload)
			})).Return(tt.mockError).Once()

			err := producer.ProduceRegistration(context.Background(), tt.user)

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
