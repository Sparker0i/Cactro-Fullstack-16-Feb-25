package entity

import "errors"

var (
	ErrInsufficientOptions   = errors.New("at least two options are required")
	ErrPollInactive          = errors.New("poll is inactive")
	ErrPollExpired           = errors.New("poll has expired")
	ErrPollNotFound          = errors.New("poll not found")
	ErrInvalidOption         = errors.New("invalid option")
	ErrDuplicateVote         = errors.New("duplicate vote")
	ErrInvalidVoteIdentifier = errors.New("invalid vote identifier")
)
