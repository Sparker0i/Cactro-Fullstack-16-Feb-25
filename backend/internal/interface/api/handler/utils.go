package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/domain/entity"
)

// Hashing utilities
func hashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(hash[:])
}

func hashFingerprint(fingerprint string) string {
	hash := sha256.Sum256([]byte(fingerprint))
	return hex.EncodeToString(hash[:])
}

// Time utilities
func isExpired(t *time.Time) bool {
	if t == nil {
		return false
	}
	return t.Before(time.Now())
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func parseTime(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %w", err)
	}
	return &t, nil
}

// String utilities
func sanitizeString(s string) string {
	return strings.TrimSpace(s)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Conversion utilities
func toJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func fromJSON(s string, v interface{}) error {
	return json.Unmarshal([]byte(s), v)
}

func toPollStatsResponse(stats *entity.PollStats) PollStatsResponse {
	options := make([]OptionResponse, len(stats.Options))
	for i, opt := range stats.Options {
		options[i] = OptionResponse{
			ID:         opt.OptionID,
			VoteCount:  opt.VoteCount,
			Percentage: roundPercentage(opt.Percentage),
		}
	}

	return PollStatsResponse{
		TotalVotes: stats.TotalVotes,
		Options:    options,
	}
}

// Math utilities
func roundPercentage(p float64) float64 {
	return float64(int(p*100+0.5)) / 100
}

// Slice utilities
func containsString(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func uniqueStrings(slice []string) []string {
	keys := make(map[string]bool)
	unique := make([]string, 0)
	for _, str := range slice {
		if !keys[str] {
			keys[str] = true
			unique = append(unique, str)
		}
	}
	return unique
}

// Error utilities
func isValidationError(err error) bool {
	_, ok := err.(ValidationErrors)
	return ok
}

func combineErrors(errs ...error) error {
	var messages []string
	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}
	if len(messages) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(messages, "; "))
}
