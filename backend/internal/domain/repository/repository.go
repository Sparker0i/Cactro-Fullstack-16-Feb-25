package repository

import (
	"context"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/google/uuid"
)

type PollRepository interface {
	Create(ctx context.Context, poll *entity.Poll) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Poll, error)
	Update(ctx context.Context, poll *entity.Poll) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, page, limit int) ([]*entity.Poll, error)
}

type VoteRepository interface {
	Create(ctx context.Context, vote *entity.Vote) error
	HasVoted(ctx context.Context, pollID uuid.UUID, identifier entity.VoteIdentifier) (bool, error)
	GetPollStats(ctx context.Context, pollID uuid.UUID) (*entity.PollStats, error)
}

type TransactionManager interface {
	Begin(ctx context.Context) (Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}
