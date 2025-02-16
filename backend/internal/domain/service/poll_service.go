package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/google/uuid"
)

type PollService interface {
	CreatePoll(ctx context.Context, question string, options []string, expiresAt *time.Time) (*entity.Poll, error)
	GetPoll(ctx context.Context, id uuid.UUID) (*entity.Poll, error)
	Vote(ctx context.Context, pollID, optionID uuid.UUID, identifier entity.VoteIdentifier) error
	ListPolls(ctx context.Context, page, limit int) ([]*entity.Poll, error)
	DeletePoll(ctx context.Context, id uuid.UUID) error
	UpdatePoll(ctx context.Context, id uuid.UUID, question string, isActive bool, expiresAt *time.Time) error
	GetPollStats(ctx context.Context, id uuid.UUID) (*entity.PollStats, error)
}

type pollService struct {
	pollRepo  repository.PollRepository
	voteRepo  repository.VoteRepository
	txManager repository.TransactionManager
	eventBus  EventBus
}

func NewPollService(
	pollRepo repository.PollRepository,
	voteRepo repository.VoteRepository,
	txManager repository.TransactionManager,
	eventBus EventBus,
) PollService {
	return &pollService{
		pollRepo:  pollRepo,
		voteRepo:  voteRepo,
		txManager: txManager,
		eventBus:  eventBus,
	}
}

func (s *pollService) CreatePoll(ctx context.Context, question string, options []string, expiresAt *time.Time) (*entity.Poll, error) {
	poll, err := entity.NewPoll(question, options, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create poll: %w", err)
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.pollRepo.Create(ctx, poll); err != nil {
		return nil, fmt.Errorf("failed to save poll: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.eventBus.Publish(PollCreatedEvent{Poll: poll})
	return poll, nil
}

func (s *pollService) GetPoll(ctx context.Context, id uuid.UUID) (*entity.Poll, error) {
	poll, err := s.pollRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}
	return poll, nil
}

func (s *pollService) Vote(ctx context.Context, pollID, optionID uuid.UUID, identifier entity.VoteIdentifier) error {
	if err := identifier.Validate(); err != nil {
		return fmt.Errorf("invalid vote identifier: %w", err)
	}

	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if already voted
	hasVoted, err := s.voteRepo.HasVoted(ctx, pollID, identifier)
	if err != nil {
		return fmt.Errorf("failed to check vote status: %w", err)
	}
	if hasVoted {
		return entity.ErrDuplicateVote
	}

	poll, err := s.pollRepo.GetByID(ctx, pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll: %w", err)
	}

	vote, err := poll.Vote(optionID, identifier)
	if err != nil {
		return fmt.Errorf("failed to record vote: %w", err)
	}

	if err := s.voteRepo.Create(ctx, vote); err != nil {
		return fmt.Errorf("failed to save vote: %w", err)
	}

	if err := s.pollRepo.Update(ctx, poll); err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	s.eventBus.Publish(VoteRecordedEvent{
		Vote: vote,
		Poll: poll,
	})

	return nil
}

func (s *pollService) ListPolls(ctx context.Context, page, limit int) ([]*entity.Poll, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	polls, err := s.pollRepo.List(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list polls: %w", err)
	}
	return polls, nil
}

func (s *pollService) DeletePoll(ctx context.Context, id uuid.UUID) error {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := s.pollRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *pollService) UpdatePoll(ctx context.Context, id uuid.UUID, question string, isActive bool, expiresAt *time.Time) error {
	tx, err := s.txManager.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	poll, err := s.pollRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get poll: %w", err)
	}

	poll.Question = question
	poll.IsActive = isActive
	poll.ExpiresAt = expiresAt
	poll.UpdatedAt = time.Now()

	if err := s.pollRepo.Update(ctx, poll); err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *pollService) GetPollStats(ctx context.Context, id uuid.UUID) (*entity.PollStats, error) {
	// First check if poll exists
	_, err := s.pollRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	stats, err := s.voteRepo.GetPollStats(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll stats: %w", err)
	}

	return stats, nil
}
