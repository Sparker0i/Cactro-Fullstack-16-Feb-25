package postgres

import (
	"context"
	"fmt"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type voteRepository struct {
	db *pgxpool.Pool
}

func NewVoteRepository(db *pgxpool.Pool) repository.VoteRepository {
	return &voteRepository{db: db}
}

func (r *voteRepository) Create(ctx context.Context, vote *entity.Vote) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO votes (id, poll_id, option_id, ip_hash, fingerprint_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		vote.ID, vote.PollID, vote.OptionID, vote.IPHash, vote.FingerprintHash, vote.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create vote: %w", err)
	}

	return nil
}

func (r *voteRepository) HasVoted(ctx context.Context, pollID uuid.UUID, identifier entity.VoteIdentifier) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM votes 
			WHERE poll_id = $1 
			AND ip_hash = $2 
			AND fingerprint_hash = $3
		)`,
		pollID, identifier.IPHash, identifier.FingerprintHash,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("failed to check vote existence: %w", err)
	}

	return exists, nil
}

func (r *voteRepository) GetPollStats(ctx context.Context, pollID uuid.UUID) (*entity.PollStats, error) {
	rows, err := r.db.Query(ctx,
		`SELECT o.id, COUNT(v.id) as vote_count
		FROM options o
		LEFT JOIN votes v ON o.id = v.option_id
		WHERE o.poll_id = $1
		GROUP BY o.id
		ORDER BY o.created_at`,
		pollID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll stats: %w", err)
	}
	defer rows.Close()

	stats := &entity.PollStats{
		Options: make([]entity.OptionStats, 0),
	}

	for rows.Next() {
		var optionStats entity.OptionStats
		err := rows.Scan(&optionStats.OptionID, &optionStats.VoteCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan option stats: %w", err)
		}
		stats.TotalVotes += optionStats.VoteCount
		stats.Options = append(stats.Options, optionStats)
	}

	// Calculate percentages
	if stats.TotalVotes > 0 {
		for i := range stats.Options {
			stats.Options[i].Percentage = float64(stats.Options[i].VoteCount) / float64(stats.TotalVotes) * 100
		}
	}

	return stats, nil
}
