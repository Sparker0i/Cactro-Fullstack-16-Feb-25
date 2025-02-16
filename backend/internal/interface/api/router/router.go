package router

import (
	"github.com/Sparker0i/cactro-polls/internal/interface/api/handler"
	"github.com/Sparker0i/cactro-polls/internal/interface/api/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine     *gin.Engine
	handler    *handler.PollHandler
	middleware *middleware.Middleware
}

func NewRouter(
	handler *handler.PollHandler,
	middleware *middleware.Middleware,
) *Router {
	return &Router{
		engine:     gin.New(),
		handler:    handler,
		middleware: middleware,
	}
}

func (r *Router) Setup() {
	// Middleware
	r.engine.Use(r.middleware.RequestID())
	r.engine.Use(r.middleware.Logger())
	r.engine.Use(r.middleware.Recovery())

	// API routes
	api := r.engine.Group("/api")
	{
		polls := api.Group("/polls")
		{
			polls.POST("", r.handler.CreatePoll)
			polls.GET("", r.handler.ListPolls)
			polls.GET("/:id", r.handler.GetPoll)
			polls.POST("/:id/vote", r.handler.Vote)
		}
	}

	// Health check
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})
}

func (r *Router) Engine() *gin.Engine {
	return r.engine
}
