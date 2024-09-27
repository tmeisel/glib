package backoff

import (
	"context"
	"errors"
	"time"

	"github.com/sethvargo/go-retry"
)

type OptionFn func(options *Backoff)

type Backoff struct {
	backoff retry.Backoff
}

var (
	ErrInvalidStrategy = errors.New("invalid Strategy specified")
	ErrStrategyNone    = errors.New("this strategy none does not allow any options")
)

type Strategy string

const (
	None        = Strategy("none")
	Constant    = Strategy("constant")
	Fibonacci   = Strategy("fibonacci")
	Exponential = Strategy("exponential")
)

func New(strategy Strategy, initial time.Duration, options ...OptionFn) (*Backoff, error) {
	var backoff retry.Backoff
	switch strategy {
	case None:
		// backoff is nil
	case Constant:
		backoff = retry.NewConstant(initial)
	case Fibonacci:
		backoff = retry.NewFibonacci(initial)
	case Exponential:
		backoff = retry.NewExponential(initial)
	default:
		return nil, ErrInvalidStrategy
	}

	r := &Backoff{
		backoff: backoff,
	}

	for _, opt := range options {
		opt(r)
	}

	return r, nil
}

func NewConstant(initial time.Duration) *Backoff {
	return &Backoff{backoff: retry.NewConstant(initial)}
}

func NewFibonacci(initial time.Duration) *Backoff {
	return &Backoff{backoff: retry.NewFibonacci(initial)}
}

func NewExponential(initial time.Duration) *Backoff {
	return &Backoff{backoff: retry.NewExponential(initial)}
}

func (b *Backoff) With(options ...OptionFn) {
	for _, opt := range options {
		opt(b)
	}
}

// Do runs the given function fn and retries on
func (b *Backoff) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if b.backoff == nil {
		return fn(ctx)
	}

	return retry.Do(ctx, b.backoff, fn)
}

// RetryableError is a wrapper to indicate a temporary
// error
var RetryableError = retry.RetryableError

// WithJitter will change the wait time randomly by
// a value between -duration and +duration
func WithJitter(duration time.Duration) OptionFn {
	return func(b *Backoff) {
		if b.backoff == nil {
			panic(ErrStrategyNone)
		}

		b.backoff = retry.WithJitter(duration, b.backoff)
	}
}

// WithCap will limit the wait time of a single retry
// to the given duration
func WithCap(duration time.Duration) OptionFn {
	return func(b *Backoff) {
		if b.backoff == nil {
			panic(ErrStrategyNone)
		}

		b.backoff = retry.WithCappedDuration(duration, b.backoff)
	}
}

// WithMaxDuration will stop retrying after the given
// duration
func WithMaxDuration(duration time.Duration) OptionFn {
	return func(b *Backoff) {
		if b.backoff == nil {
			panic(ErrStrategyNone)
		}

		b.backoff = retry.WithMaxDuration(duration, b.backoff)
	}
}

// WithMaxRetries limits the number of retries to n
func WithMaxRetries(n uint64) OptionFn {
	return func(b *Backoff) {
		if b.backoff == nil {
			panic(ErrStrategyNone)
		}

		b.backoff = retry.WithMaxRetries(n, b.backoff)
	}
}
