package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test cases
func TestPollService_CreatePoll(t *testing.T) {
	ctx := context.Background()
	pollRepo := new(MockPollRepository)
	voteRepo := new(MockVoteRepository)
	txManager := new(MockTransactionManager)
	eventBus := new(MockEventBus)
	tx := new(MockTransaction)

	pollService := service.NewPollService(pollRepo, voteRepo, txManager, eventBus)

	tests := []struct {
		name      string
		question  string
		options   []string
		expiresAt *time.Time
		mockSetup func()
		wantErr   bool
	}{
		{
			name:     "Successful poll creation",
			question: "Test question?",
			options:  []string{"Option 1", "Option 2"},
			mockSetup: func() {
				txManager.On("Begin", ctx).Return(tx, nil)
				pollRepo.On("Create", ctx, mock.AnythingOfType("*entity.Poll")).Return(nil)
				tx.On("Commit").Return(nil)
				eventBus.On("Publish", mock.AnythingOfType("service.PollCreatedEvent")).Return()
			},
			wantErr: false,
		},
		{
			name:     "Failed poll creation - database error",
			question: "Test question?",
			options:  []string{"Option 1", "Option 2"},
			mockSetup: func() {
				txManager.On("Begin", ctx).Return(tx, nil)
				pollRepo.On("Create", ctx, mock.AnythingOfType("*entity.Poll")).Return(assert.AnError)
				tx.On("Rollback").Return(nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			tt.mockSetup()

			// Execute test
			poll, err := pollService.CreatePoll(ctx, tt.question, tt.options, tt.expiresAt)

			// Assert results
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, poll)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, poll)
				assert.Equal(t, tt.question, poll.Question)
				assert.Len(t, poll.Options, len(tt.options))
			}

			// Verify mock expectations
			pollRepo.AssertExpectations(t)
			txManager.AssertExpectations(t)
			tx.AssertExpectations(t)
			eventBus.AssertExpectations(t)
		})
	}
}
