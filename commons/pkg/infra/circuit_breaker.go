package infra

import (
	"errors"
	"time"
)

type CircuitBreaker[T any] struct {
	maxFailures  int
	resetTimeout time.Duration
	failures     int
	lastFailure  time.Time
}

var ErrorCircuitBreakerOpen = errors.New("circuit breaker open")

func NewCircuitBreaker[T any](maxFailures int, resetTimeout time.Duration) *CircuitBreaker[T] {
	return &CircuitBreaker[T]{
		maxFailures:  maxFailures,
		resetTimeout: resetTimeout,
		failures:     0,
		lastFailure:  time.Time{},
	}
}

func (cb *CircuitBreaker[T]) Call(f func() (T, error)) (T, error) {
	tryAgain := false
	var result T
	var err error

	if time.Since(cb.lastFailure) < cb.resetTimeout {
		return result, ErrorCircuitBreakerOpen
	}

	for tryAgain {
		tryAgain = false
		result, err = f()

		if err != nil {
			cb.failures++
			if cb.failures >= cb.maxFailures {
				cb.lastFailure = time.Now()
			}
			tryAgain = true
		} else {
			cb.failures = 0
		}
	}

	return result, err
}
