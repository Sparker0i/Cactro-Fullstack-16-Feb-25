package container

import (
	"sync"

	"github.com/Sparker0i/cactro-polls/internal/domain/repository"
	"github.com/Sparker0i/cactro-polls/internal/domain/service"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/config"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/database"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/event"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/logger"
	"github.com/Sparker0i/cactro-polls/internal/infrastructure/server"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/handler"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/middleware"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/router"
	"github.com/Sparker0i/cactro-polls/internal/interface/repository/postgres"
)

type Container struct {
	cfg    *config.Config
	mu     sync.Mutex
	cached map[string]interface{}
}

func NewContainer(cfg *config.Config) *Container {
	return &Container{
		cfg:    cfg,
		cached: make(map[string]interface{}),
	}
}

func (c *Container) GetLogger() logger.Logger {
	return c.singleton("logger", func() interface{} {
		log, err := logger.NewLogger(&c.cfg.Logger)
		if err != nil {
			panic(err)
		}
		return log
	}).(logger.Logger)
}

func (c *Container) GetDatabase() *database.Database {
	return c.singleton("database", func() interface{} {
		db, err := database.NewDatabase(&c.cfg.Database)
		if err != nil {
			panic(err)
		}
		return db
	}).(*database.Database)
}

func (c *Container) GetEventBus() event.EventBus {
	return c.singleton("event_bus", func() interface{} {
		return event.NewEventBus()
	}).(event.EventBus)
}

func (c *Container) GetPollRepository() repository.PollRepository {
	return c.singleton("poll_repository", func() interface{} {
		return postgres.NewPollRepository(c.GetDatabase().Pool())
	}).(repository.PollRepository)
}

func (c *Container) GetVoteRepository() repository.VoteRepository {
	return c.singleton("vote_repository", func() interface{} {
		return postgres.NewVoteRepository(c.GetDatabase().Pool())
	}).(repository.VoteRepository)
}

func (c *Container) GetTransactionManager() repository.TransactionManager {
	return c.singleton("transaction_manager", func() interface{} {
		return postgres.NewTransactionManager(c.GetDatabase().Pool())
	}).(repository.TransactionManager)
}

func (c *Container) GetPollService() service.PollService {
	return c.singleton("poll_service", func() interface{} {
		return service.NewPollService(
			c.GetPollRepository(),
			c.GetVoteRepository(),
			c.GetTransactionManager(),
			c.GetEventBus(),
		)
	}).(service.PollService)
}

func (c *Container) GetMiddleware() *middleware.Middleware {
	return c.singleton("middleware", func() interface{} {
		return middleware.NewMiddleware(c.GetLogger())
	}).(*middleware.Middleware)
}

func (c *Container) GetPollHandler() *handler.PollHandler {
	return c.singleton("poll_handler", func() interface{} {
		return handler.NewPollHandler(c.GetPollService())
	}).(*handler.PollHandler)
}

func (c *Container) GetRouter() *router.Router {
	return c.singleton("router", func() interface{} {
		r := router.NewRouter(
			c.GetPollHandler(),
			c.GetMiddleware(),
		)
		r.Setup()
		return r
	}).(*router.Router)
}

func (c *Container) GetHTTPServer() *server.Server {
	return c.singleton("http_server", func() interface{} {
		return server.NewServer(
			c.GetRouter(),
			c.GetLogger(),
			c.cfg,
		)
	}).(*server.Server)
}

func (c *Container) singleton(key string, factory func() interface{}) interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	if instance, exists := c.cached[key]; exists {
		return instance
	}

	instance := factory()
	c.cached[key] = instance
	return instance
}

func (c *Container) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Close database connection
	if db, ok := c.cached["database"].(*database.Database); ok {
		db.Close()
	}

	// Stop event bus
	if eventBus, ok := c.cached["event_bus"].(event.EventBus); ok {
		eventBus.Stop()
	}

	// Clear cache
	c.cached = make(map[string]interface{})
}
