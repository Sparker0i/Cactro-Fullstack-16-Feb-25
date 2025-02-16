package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

// Validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()

	// Register custom validators
	validate.RegisterValidation("future_time", validateFutureTime)
	validate.RegisterValidation("valid_option", validateOption)
	validate.RegisterValidation("valid_question", validateQuestion)
}

// Custom validators
func validateFutureTime(fl validator.FieldLevel) bool {
	timeVal, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	return timeVal.After(time.Now())
}

func validateOption(fl validator.FieldLevel) bool {
	option := fl.Field().String()
	trimmed := strings.TrimSpace(option)
	return len(trimmed) >= 1 && len(trimmed) <= 200
}

func validateQuestion(fl validator.FieldLevel) bool {
	question := fl.Field().String()
	trimmed := strings.TrimSpace(question)
	return len(trimmed) >= 5 && len(trimmed) <= 500
}

// Validation helpers
func validateRequest(c *gin.Context, req interface{}) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return fmt.Errorf("invalid request format: %w", err)
	}

	if err := validate.Struct(req); err != nil {
		return formatValidationErrors(err)
	}

	return nil
}

func validatePollID(id string) (uuid.UUID, error) {
	pollID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid poll ID format")
	}
	return pollID, nil
}

func validatePaginationParams(page, limit string) (int, int) {
	pageNum := 1
	if p, err := strconv.Atoi(page); err == nil && p > 0 {
		pageNum = p
	}

	limitNum := 10
	if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
		limitNum = l
	}

	return pageNum, limitNum
}

func (ve ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Error formatting
func formatValidationErrors(err error) error {
	var errors ValidationErrors

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fmt.Errorf("invalid request format")
	}

	for _, e := range validationErrors {
		errors = append(errors, ValidationError{
			Field:   e.Field(),
			Message: getValidationErrorMessage(e),
		})
	}

	return errors
}

func getValidationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return fmt.Sprintf("Minimum length is %s", e.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", e.Param())
	case "future_time":
		return "Time must be in the future"
	case "valid_option":
		return "Option must be between 1 and 200 characters"
	case "valid_question":
		return "Question must be between 5 and 500 characters"
	default:
		return fmt.Sprintf("Validation failed on condition: %s", e.Tag())
	}
}
