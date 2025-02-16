package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type pollRepository struct {
	db *pgxpool.Pool
}

func NewPollRepository(db *pgxpool.Pool) repository.PollRepository {
	return &pollRepository{db: db}
}

func (r *pollRepository) Create(ctx context.Context, poll *entity.Poll) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Insert poll
	_, err = tx.Exec(ctx,
		`INSERT INTO polls (id, question, expires_at, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		poll.ID, poll.Question, poll.ExpiresAt, poll.IsActive, poll.CreatedAt, poll.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert poll: %w", err)
	}

	// Insert options
	for _, option := range poll.Options {
		_, err = tx.Exec(ctx,
			`INSERT INTO options (id, poll_id, option_text, created_at)
			VALUES ($1, $2, $3, $4)`,
			option.ID, poll.ID, option.OptionText, option.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert option: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *pollRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Poll, error) {
	var poll entity.Poll

	err := r.db.QueryRow(ctx,
		`SELECT id, question, expires_at, is_active, created_at, updated_at
		FROM polls WHERE id = $1`,
		id,
	).Scan(
		&poll.ID,
		&poll.Question,
		&poll.ExpiresAt,
		&poll.IsActive,
		&poll.CreatedAt,
		&poll.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrPollNotFound
		}
		return nil, fmt.Errorf("failed to get poll: %w", err)
	}

	// Get options with vote counts
	rows, err := r.db.Query(ctx,
		`SELECT o.id, o.option_text, o.created_at, COUNT(v.id) as vote_count
		FROM options o
		LEFT JOIN votes v ON o.id = v.option_id
		WHERE o.poll_id = $1
		GROUP BY o.id, o.option_text, o.created_at
		ORDER BY o.created_at`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll options: %w", err)
	}
	defer rows.Close()

	poll.Options = make([]entity.Option, 0)
	var totalVotes int

	for rows.Next() {
		var option entity.Option
		err := rows.Scan(
			&option.ID,
			&option.OptionText,
			&option.CreatedAt,
			&option.VoteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan option: %w", err)
		}
		option.PollID = poll.ID
		totalVotes += option.VoteCount
		poll.Options = append(poll.Options, option)
	}

	// Calculate percentages
	if totalVotes > 0 {
		for i := range poll.Options {
			poll.Options[i].Percentage = float64(poll.Options[i].VoteCount) / float64(totalVotes) * 100
		}
	}

	return &poll, nil
}

func (r *pollRepository) Update(ctx context.Context, poll *entity.Poll) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`UPDATE polls 
		SET question = $1, expires_at = $2, is_active = $3, updated_at = $4
		WHERE id = $5`,
		poll.Question, poll.ExpiresAt, poll.IsActive, poll.UpdatedAt, poll.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update poll: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *pollRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		`DELETE FROM polls WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete poll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return entity.ErrPollNotFound
	}

	return nil
}

func (r *pollRepository) List(ctx context.Context, page, limit int) ([]*entity.Poll, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Query(ctx,
		`SELECT id, question, expires_at, is_active, created_at, updated_at
		FROM polls
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list polls: %w", err)
	}
	defer rows.Close()

	polls := make([]*entity.Poll, 0)
	for rows.Next() {
		var poll entity.Poll
		err := rows.Scan(
			&poll.ID,
			&poll.Question,
			&poll.ExpiresAt,
			&poll.IsActive,
			&poll.CreatedAt,
			&poll.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan poll: %w", err)
		}
		polls = append(polls, &poll)
	}

	return polls, nil
}
