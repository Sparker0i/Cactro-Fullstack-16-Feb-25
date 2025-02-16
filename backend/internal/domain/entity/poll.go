package entity

import (
	"time"

	"github.com/google/uuid"
)

type Poll struct {
	ID        uuid.UUID
	Question  string
	Options   []Option
	CreatedAt time.Time
	ExpiresAt *time.Time
	IsActive  bool
	UpdatedAt time.Time
}

type Option struct {
	ID         uuid.UUID
	PollID     uuid.UUID
	OptionText string
	CreatedAt  time.Time
	VoteCount  int
	Percentage float64
}

type Vote struct {
	ID              uuid.UUID
	PollID          uuid.UUID
	OptionID        uuid.UUID
	IPHash          string
	FingerprintHash string
	CreatedAt       time.Time
}

type PollStats struct {
	TotalVotes int
	Options    []OptionStats
}

type OptionStats struct {
	OptionID   uuid.UUID
	VoteCount  int
	Percentage float64
}

// NewPoll creates a new poll with the given options
func NewPoll(question string, options []string, expiresAt *time.Time) (*Poll, error) {
	if len(options) < 2 {
		return nil, ErrInsufficientOptions
	}

	pollID := uuid.New()
	now := time.Now()

	pollOptions := make([]Option, len(options))
	for i, opt := range options {
		pollOptions[i] = Option{
			ID:         uuid.New(),
			PollID:     pollID,
			OptionText: opt,
			CreatedAt:  now,
		}
	}

	return &Poll{
		ID:        pollID,
		Question:  question,
		Options:   pollOptions,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		IsActive:  true,
		UpdatedAt: now,
	}, nil
}

// Vote records a vote for the given option
func (p *Poll) Vote(optionID uuid.UUID, identifier VoteIdentifier) (*Vote, error) {
	if !p.IsActive {
		return nil, ErrPollInactive
	}

	if p.ExpiresAt != nil && p.ExpiresAt.Before(time.Now()) {
		return nil, ErrPollExpired
	}

	var targetOption *Option
	for i := range p.Options {
		if p.Options[i].ID == optionID {
			targetOption = &p.Options[i]
			break
		}
	}

	if targetOption == nil {
		return nil, ErrInvalidOption
	}

	vote := &Vote{
		ID:              uuid.New(),
		PollID:          p.ID,
		OptionID:        optionID,
		IPHash:          identifier.IPHash,
		FingerprintHash: identifier.FingerprintHash,
		CreatedAt:       time.Now(),
	}

	targetOption.VoteCount++
	p.updatePercentages()

	return vote, nil
}

func (p *Poll) updatePercentages() {
	total := 0
	for _, opt := range p.Options {
		total += opt.VoteCount
	}

	if total > 0 {
		for i := range p.Options {
			p.Options[i].Percentage = float64(p.Options[i].VoteCount) / float64(total) * 100
		}
	}
}
