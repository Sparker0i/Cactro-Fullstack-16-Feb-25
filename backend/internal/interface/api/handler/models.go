package handler

import (
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/google/uuid"
)

// Request models
type CreatePollRequest struct {
	Question  string    `json:"question" binding:"required,min=5,max=500"`
	Options   []string  `json:"options" binding:"required,min=2,dive,required"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type UpdatePollRequest struct {
	Question  string     `json:"question" binding:"omitempty,min=5,max=500"`
	IsActive  *bool      `json:"is_active,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type VoteRequest struct {
	OptionID        uuid.UUID `json:"option_id" binding:"required"`
	FingerprintHash string    `json:"fingerprint_hash" binding:"required,min=32"`
}

// Response models
type PollResponse struct {
	ID        uuid.UUID        `json:"id"`
	Question  string           `json:"question"`
	Options   []OptionResponse `json:"options"`
	CreatedAt time.Time        `json:"created_at"`
	ExpiresAt *time.Time       `json:"expires_at,omitempty"`
	IsActive  bool             `json:"is_active"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type OptionResponse struct {
	ID         uuid.UUID `json:"id"`
	OptionText string    `json:"option_text"`
	VoteCount  int       `json:"vote_count"`
	Percentage float64   `json:"percentage"`
}

type PollStatsResponse struct {
	TotalVotes int              `json:"total_votes"`
	Options    []OptionResponse `json:"options"`
}

type PollListResponse struct {
	Polls      []PollResponse `json:"polls"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPolls int            `json:"total_polls"`
}

// Converters
func toPollResponse(poll *entity.Poll) PollResponse {
	options := make([]OptionResponse, len(poll.Options))
	for i, opt := range poll.Options {
		options[i] = OptionResponse{
			ID:         opt.ID,
			OptionText: opt.OptionText,
			VoteCount:  opt.VoteCount,
			Percentage: opt.Percentage,
		}
	}

	return PollResponse{
		ID:        poll.ID,
		Question:  poll.Question,
		Options:   options,
		CreatedAt: poll.CreatedAt,
		ExpiresAt: poll.ExpiresAt,
		IsActive:  poll.IsActive,
		UpdatedAt: poll.UpdatedAt,
	}
}
