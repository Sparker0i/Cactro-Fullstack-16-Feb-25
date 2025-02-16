package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
	"github.com/Sparker0i/cactro-polls/internal/domain/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PollHandler struct {
	pollService service.PollService
}

func NewPollHandler(pollService service.PollService) *PollHandler {
	return &PollHandler{
		pollService: pollService,
	}
}

// CreatePoll godoc
// @Summary Create a new poll
// @Description Create a new poll with options
// @Tags polls
// @Accept json
// @Produce json
// @Param poll body CreatePollRequest true "Poll to create"
// @Success 201 {object} PollResponse
// @Failure 400 {object} ErrorResponse
// @Router /polls [post]
func (h *PollHandler) CreatePoll(c *gin.Context) {
	var req CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	poll, err := h.pollService.CreatePoll(c.Request.Context(), req.Question, req.Options, &req.ExpiresAt)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, toPollResponse(poll))
}

// GetPoll godoc
// @Summary Get a poll by ID
// @Description Get a poll's details including options and vote counts
// @Tags polls
// @Produce json
// @Param id path string true "Poll ID"
// @Success 200 {object} PollResponse
// @Failure 404 {object} ErrorResponse
// @Router /polls/{id} [get]
func (h *PollHandler) GetPoll(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	poll, err := h.pollService.GetPoll(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPollResponse(poll))
}

// Vote godoc
// @Summary Cast a vote for a poll option
// @Description Cast a vote for a specific option in a poll
// @Tags polls
// @Accept json
// @Produce json
// @Param id path string true "Poll ID"
// @Param vote body VoteRequest true "Vote details"
// @Success 200 {object} PollStatsResponse
// @Failure 400,404,409 {object} ErrorResponse
// @Router /polls/{id}/vote [post]
func (h *PollHandler) Vote(c *gin.Context) {
	pollID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	var req VoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err)
		return
	}

	identifier := entity.VoteIdentifier{
		IPHash:          hashIP(c.ClientIP()),
		FingerprintHash: req.FingerprintHash,
	}

	err = h.pollService.Vote(c.Request.Context(), pollID, req.OptionID, identifier)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	// Get updated stats
	stats, err := h.pollService.GetPollStats(c.Request.Context(), pollID)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPollStatsResponse(stats))
}

// ListPolls godoc
// @Summary List all polls
// @Description Get a paginated list of polls
// @Tags polls
// @Produce json
// @Param page query integer false "Page number"
// @Param limit query integer false "Items per page"
// @Success 200 {object} PollListResponse
// @Router /polls [get]
func (h *PollHandler) ListPolls(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	polls, err := h.pollService.ListPolls(c.Request.Context(), page, limit)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response := PollListResponse{
		Polls:    make([]PollResponse, len(polls)),
		Page:     page,
		PageSize: limit,
	}

	for i, poll := range polls {
		response.Polls[i] = toPollResponse(poll)
	}

	c.JSON(http.StatusOK, response)
}

// Error handling helpers
func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrPollNotFound):
		respondWithError(c, http.StatusNotFound, err)
	case errors.Is(err, entity.ErrDuplicateVote):
		respondWithError(c, http.StatusConflict, err)
	case errors.Is(err, entity.ErrPollInactive):
		respondWithError(c, http.StatusForbidden, err)
	case errors.Is(err, entity.ErrPollExpired):
		respondWithError(c, http.StatusForbidden, err)
	default:
		respondWithError(c, http.StatusInternalServerError, err)
	}
}
