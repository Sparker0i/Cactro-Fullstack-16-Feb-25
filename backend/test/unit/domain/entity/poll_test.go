package entity_test

import (
	"testing"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestNewPoll(t *testing.T) {
	tests := []struct {
		name      string
		question  string
		options   []string
		expiresAt *time.Time
		wantErr   error
	}{
		{
			name:     "Valid poll",
			question: "Test question?",
			options:  []string{"Option 1", "Option 2"},
			wantErr:  nil,
		},
		{
			name:     "Invalid poll - insufficient options",
			question: "Test question?",
			options:  []string{"Option 1"},
			wantErr:  entity.ErrInsufficientOptions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poll, err := entity.NewPoll(tt.question, tt.options, tt.expiresAt)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, poll)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, poll)
			assert.Equal(t, tt.question, poll.Question)
			assert.Len(t, poll.Options, len(tt.options))
			assert.True(t, poll.IsActive)
		})
	}
}

func TestPoll_Vote(t *testing.T) {
	now := time.Now()
	expiredTime := now.Add(-24 * time.Hour)
	futureTime := now.Add(24 * time.Hour)

	tests := []struct {
		name       string
		setupPoll  func() *entity.Poll
		identifier entity.VoteIdentifier
		wantErr    error
	}{
		{
			name: "Valid vote",
			setupPoll: func() *entity.Poll {
				poll, _ := entity.NewPoll("Test?", []string{"A", "B"}, &futureTime)
				return poll
			},
			identifier: entity.VoteIdentifier{
				IPHash:          "testhash",
				FingerprintHash: "fingerprintHash",
			},
			wantErr: nil,
		},
		{
			name: "Vote on expired poll",
			setupPoll: func() *entity.Poll {
				poll, _ := entity.NewPoll("Test?", []string{"A", "B"}, &expiredTime)
				return poll
			},
			identifier: entity.VoteIdentifier{
				IPHash:          "testhash",
				FingerprintHash: "fingerprintHash",
			},
			wantErr: entity.ErrPollExpired,
		},
		{
			name: "Vote on inactive poll",
			setupPoll: func() *entity.Poll {
				poll, _ := entity.NewPoll("Test?", []string{"A", "B"}, nil)
				poll.IsActive = false
				return poll
			},
			identifier: entity.VoteIdentifier{
				IPHash:          "testhash",
				FingerprintHash: "fingerprintHash",
			},
			wantErr: entity.ErrPollInactive,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poll := tt.setupPoll()
			optionID := poll.Options[0].ID

			vote, err := poll.Vote(optionID, tt.identifier)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, vote)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, vote)
			assert.Equal(t, poll.ID, vote.PollID)
			assert.Equal(t, optionID, vote.OptionID)
			assert.Equal(t, tt.identifier.IPHash, vote.IPHash)
			assert.Equal(t, tt.identifier.FingerprintHash, vote.FingerprintHash)
		})
	}
}
