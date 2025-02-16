package service_test

import (
	"context"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPollRepository implements repository.PollRepository
type MockPollRepository struct {
	mock.Mock
}

func (m *MockPollRepository) Create(ctx context.Context, poll *entity.Poll) error {
	args := m.Called(ctx, poll)
	return args.Error(0)
}

func (m *MockPollRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Poll, error) {
	args := m.Called(ctx, id)
	if poll, ok := args.Get(0).(*entity.Poll); ok {
		return poll, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPollRepository) Update(ctx context.Context, poll *entity.Poll) error {
	args := m.Called(ctx, poll)
	return args.Error(0)
}

func (m *MockPollRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPollRepository) List(ctx context.Context, page, limit int) ([]*entity.Poll, error) {
	args := m.Called(ctx, page, limit)
	if polls, ok := args.Get(0).([]*entity.Poll); ok {
		return polls, args.Error(1)
	}
	return nil, args.Error(1)
}

// MockVoteRepository implements repository.VoteRepository
type MockVoteRepository struct {
	mock.Mock
}

func (m *MockVoteRepository) Create(ctx context.Context, vote *entity.Vote) error {
	args := m.Called(ctx, vote)
	return args.Error(0)
}

func (m *MockVoteRepository) HasVoted(ctx context.Context, pollID uuid.UUID, identifier entity.VoteIdentifier) (bool, error) {
	args := m.Called(ctx, pollID, identifier)
	return args.Bool(0), args.Error(1)
}

func (m *MockVoteRepository) GetPollStats(ctx context.Context, pollID uuid.UUID) (*entity.PollStats, error) {
	args := m.Called(ctx, pollID)
	if stats, ok := args.Get(0).(*entity.PollStats); ok {
		return stats, args.Error(1)
	}
	return nil, args.Error(1)
}

// MockTransactionManager implements repository.TransactionManager
type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Begin(ctx context.Context) (repository.Transaction, error) {
	args := m.Called(ctx)
	if tx, ok := args.Get(0).(repository.Transaction); ok {
		return tx, args.Error(1)
	}
	return nil, args.Error(1)
}

// MockTransaction implements repository.Transaction
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}

// MockEventBus implements service.EventBus
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(event interface{}) {
	m.Called(event)
}

func (m *MockEventBus) Stop() {
	m.Called()
}
