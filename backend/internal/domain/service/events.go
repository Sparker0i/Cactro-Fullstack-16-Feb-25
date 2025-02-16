package service

import "github.com/Sparker0i/cactro-polls/internal/domain/entity"

// EventBus defines the interface for event publishing
type EventBus interface {
	Publish(event interface{})
}

// Domain events
type PollCreatedEvent struct {
	Poll *entity.Poll
}

type VoteRecordedEvent struct {
	Vote *entity.Vote
	Poll *entity.Poll
}
