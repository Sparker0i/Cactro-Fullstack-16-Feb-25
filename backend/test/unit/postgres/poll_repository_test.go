package postgres_test

import (
	"context"
	"testing"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/database"
	"github.com/Sparker0i/cactro-polls/internal/interface/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type PollRepositoryTestSuite struct {
	suite.Suite
	db       *database.Database
	pollRepo repository.PollRepository
	ctx      context.Context
}

func TestPollRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PollRepositoryTestSuite))
}

func (s *PollRepositoryTestSuite) SetupSuite() {
	cfg := &config.DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "polling_app_test",
		SSLMode:  "disable",
	}

	db, err := database.NewDatabase(cfg)
	s.Require().NoError(err)
	s.db = db

	s.pollRepo = postgres.NewPollRepository(db.Pool())
	s.ctx = context.Background()

	// Run migrations
	// You would typically use a migration tool here
}

func (s *PollRepositoryTestSuite) TearDownSuite() {
	if s.db != nil {
		s.db.Close()
	}
}

func (s *PollRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	_, err := s.db.Pool().Exec(s.ctx, "TRUNCATE polls CASCADE")
	s.Require().NoError(err)
}

func (s *PollRepositoryTestSuite) TestCreatePoll() {
	// Create test poll
	poll, err := entity.NewPoll("Test question?", []string{"Option 1", "Option 2"}, nil)
	s.Require().NoError(err)

	// Test creation
	err = s.pollRepo.Create(s.ctx, poll)
	s.Require().NoError(err)

	// Verify creation
	saved, err := s.pollRepo.GetByID(s.ctx, poll.ID)
	s.Require().NoError(err)
	s.NotNil(saved)
	s.Equal(poll.Question, saved.Question)
	s.Len(saved.Options, len(poll.Options))
}

func (s *PollRepositoryTestSuite) TestGetByID() {
	// Create test poll
	poll, err := entity.NewPoll("Test question?", []string{"Option 1", "Option 2"}, nil)
	s.Require().NoError(err)
	err = s.pollRepo.Create(s.ctx, poll)
	s.Require().NoError(err)

	tests := []struct {
		name    string
		id      uuid.UUID
		wantErr bool
	}{
		{
			name:    "Existing poll",
			id:      poll.ID,
			wantErr: false,
		},
		{
			name:    "Non-existent poll",
			id:      uuid.New(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			saved, err := s.pollRepo.GetByID(s.ctx, tt.id)
			if tt.wantErr {
				s.Error(err)
				s.Nil(saved)
			} else {
				s.NoError(err)
				s.NotNil(saved)
				s.Equal(poll.Question, saved.Question)
			}
		})
	}
}
