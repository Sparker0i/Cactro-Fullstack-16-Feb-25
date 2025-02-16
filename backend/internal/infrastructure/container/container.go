package container

import (
	"sync"

	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/Sparker0i/cactro-polls/internal/domain/service"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/database"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/event"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/logger"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/handler"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/middleware"
	"github.com/Sparker0i/cactro-polls/internal/interface/repository/postgres"
	"github.com/gin-gonic/gin"
)

type Container struct {
	cfg        *config.Config
	mu         sync.Mutex
	logger     logger.Logger
	db         *database.Database
	engine     *gin.Engine
	components componentContainer
}

type componentContainer struct {
	eventBus    event.EventBus
	pollRepo    repository.PollRepository
	voteRepo    repository.VoteRepository
	txManager   repository.TransactionManager
	pollService service.PollService
	middleware  *middleware.Middleware
	pollHandler *handler.PollHandler
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{
		cfg: cfg,
	}

	// Initialize core components first
	if err := c.initializeCore(); err != nil {
		return nil, err
	}

	// Initialize other components
	if err := c.initializeComponents(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Container) initializeCore() error {
	// Initialize logger first
	log, err := logger.NewLogger(&c.cfg.Logger)
	if err != nil {
		return err
	}
	c.logger = log

	// Initialize database
	db, err := database.NewDatabase(&c.cfg.Database)
	if err != nil {
		return err
	}
	c.db = db

	return nil
}

func (c *Container) initializeComponents() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Initialize event bus
	c.components.eventBus = event.NewEventBus()

	// Initialize repositories
	c.components.pollRepo = postgres.NewPollRepository(c.db.Pool())
	c.components.voteRepo = postgres.NewVoteRepository(c.db.Pool())
	c.components.txManager = postgres.NewTransactionManager(c.db.Pool())

	// Initialize service
	c.components.pollService = service.NewPollService(
		c.components.pollRepo,
		c.components.voteRepo,
		c.components.txManager,
		c.components.eventBus,
	)

	// Initialize API components
	c.components.middleware = middleware.NewMiddleware(c.logger)
	c.components.pollHandler = handler.NewPollHandler(c.components.pollService)

	return nil
}

func (c *Container) InitializeHTTP() *gin.Engine {
	gin.SetMode(c.cfg.Server.Mode)
	engine := gin.New()

	// Setup middleware
	engine.Use(c.components.middleware.RequestID())
	engine.Use(c.components.middleware.Logger())
	engine.Use(c.components.middleware.Recovery())

	// Setup routes
	api := engine.Group("/api")
	{
		polls := api.Group("/polls")
		{
			polls.POST("", c.components.pollHandler.CreatePoll)
			polls.GET("", c.components.pollHandler.ListPolls)
			polls.GET("/:id", c.components.pollHandler.GetPoll)
			polls.POST("/:id/vote", c.components.pollHandler.Vote)
		}
	}

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	c.engine = engine
	return engine
}

func (c *Container) Logger() logger.Logger {
	return c.logger
}

func (c *Container) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.components.eventBus != nil {
		c.components.eventBus.Stop()
	}

	if c.db != nil {
		c.db.Close()
	}
}
