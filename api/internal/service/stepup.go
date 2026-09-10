package service

import (
	"context"
	"errors"
	"log"
	"time"
)

var (
	ErrSaveStepUp = errors.New("Unable to save verification. Please try again.")
)

// StepUpExpiration is a safety-net expiry for a grant that is created but never consumed; normal
// use consumes the grant on the very next protected request.
const StepUpExpiration = 5 * time.Minute

func stepUpCacheKey(token string) string {
	return "stepup_" + token
}

func (s *Service) SetStepUp(ctx context.Context, token string) error {
	err := s.Cache.Set(ctx, stepUpCacheKey(token), "1", StepUpExpiration)
	if err != nil {
		log.Printf("error saving step-up verification: %s", err.Error())
		return ErrSaveStepUp
	}

	return nil
}

// ConsumeStepUp is single-use: a cache miss is the expected, common case (no prior verification),
// not an exceptional one, so it is not logged like other cache lookups in this package.
func (s *Service) ConsumeStepUp(ctx context.Context, token string) (bool, error) {
	val, err := s.Cache.Get(ctx, stepUpCacheKey(token))
	if err != nil || val != "1" {
		return false, nil
	}

	_ = s.Cache.Del(ctx, stepUpCacheKey(token))

	return true, nil
}
