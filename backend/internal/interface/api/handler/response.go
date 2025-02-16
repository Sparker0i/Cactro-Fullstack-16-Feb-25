package handler

import (
	"net/http"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

// Response wrappers
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

type ErrorData struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

type MetaData struct {
	Page       int       `json:"page,omitempty"`
	PageSize   int       `json:"page_size,omitempty"`
	TotalItems int       `json:"total_items,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// Response helpers
func respondWithSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &MetaData{
			Timestamp: time.Now(),
		},
	})
}

func respondWithCreated(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
		Meta: &MetaData{
			Timestamp: time.Now(),
		},
	})
}

func respondWithError(c *gin.Context, code int, err error) {
	requestID, _ := c.Get("request_id")
	var errorCode string
	var message string

	switch {
	case err == entity.ErrPollNotFound:
		errorCode = "POLL_NOT_FOUND"
		message = "Poll not found"
	case err == entity.ErrDuplicateVote:
		errorCode = "DUPLICATE_VOTE"
		message = "You have already voted in this poll"
	case err == entity.ErrPollInactive:
		errorCode = "POLL_INACTIVE"
		message = "This poll is no longer active"
	case err == entity.ErrPollExpired:
		errorCode = "POLL_EXPIRED"
		message = "This poll has expired"
	default:
		errorCode = "INTERNAL_ERROR"
		message = "An internal error occurred"
	}

	c.JSON(code, Response{
		Success: false,
		Error: &ErrorData{
			Code:      errorCode,
			Message:   message,
			Details:   err.Error(),
			RequestID: requestID.(string),
			Timestamp: time.Now(),
		},
	})
}

func respondWithPagination(c *gin.Context, data interface{}, page, pageSize, totalItems int) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
		Meta: &MetaData{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: totalItems,
			Timestamp:  time.Now(),
		},
	})
}
