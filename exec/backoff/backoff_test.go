package backoff

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	retryableError = RetryableError(errors.New("retry"))
	fatalError     = errors.New("fatal")
)

func TestNew(t *testing.T) {
	const maxRetries = 0

	type testCase struct {
		Strategy Strategy
		Initial  time.Duration
		Options  []OptionFn
		Expected error
	}

	for name, tc := range map[string]testCase{
		"invalid Strategy": {
			Strategy: "something stupid",
			Initial:  time.Millisecond,
			Expected: ErrInvalidStrategy,
		},
		"Strategy none": {
			Strategy: None,
			Initial:  time.Millisecond,
		},
		"Strategy constant": {
			Strategy: Constant,
			Initial:  time.Millisecond,
			Options:  []OptionFn{WithMaxRetries(maxRetries)},
		},
		"Strategy fibonacci": {
			Strategy: Fibonacci,
			Initial:  time.Millisecond,
			Options:  []OptionFn{WithMaxRetries(maxRetries)},
		},
		"Strategy exponential": {
			Strategy: Exponential,
			Initial:  time.Millisecond,
			Options:  []OptionFn{WithMaxRetries(maxRetries)},
		},
	} {
		t.Run(name, func(t *testing.T) {
			b, err := New(tc.Strategy, tc.Initial, tc.Options...)

			if tc.Expected != nil {
				require.Error(t, err)
				assert.Equal(t, tc.Expected, err)
				assert.Nil(t, b)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, b)
		})
	}
}

func TestDo(t *testing.T) {
	type testCase struct {
		ExpectedTries uint64
		Errors        []error
		ExpectedError error
	}

	for name, tc := range map[string]testCase{
		"no error, no retry": {
			Errors:        []error{nil},
			ExpectedError: nil,
			ExpectedTries: 1,
		},
		"one retryable error, one retry": {
			Errors:        []error{retryableError, nil},
			ExpectedError: nil,
			ExpectedTries: 2,
		},
		"one retryable error, one fatal error": {
			Errors:        []error{retryableError, fatalError},
			ExpectedError: fatalError,
			ExpectedTries: 2,
		},
	} {
		t.Run(name, func(t *testing.T) {
			b, err := New(Constant, time.Millisecond)
			require.NoError(t, err)
			require.NotNil(t, b)

			retries := make([]time.Duration, 0)
			lastStart := time.Now()

			err = b.Do(context.Background(), func(ctx context.Context) error {
				current := len(retries)
				retries = append(retries, time.Now().Sub(lastStart))

				lastStart = time.Now()

				return tc.Errors[current]
			})

			assert.Equal(t, tc.ExpectedError, err)
			assert.Equal(t, tc.ExpectedTries, uint64(len(retries)))
		})
	}
}

func TestWithMaxRetries(t *testing.T) {
	const maxRetries = 2

	b, err := New(Constant, time.Millisecond, WithMaxRetries(maxRetries))
	require.NoError(t, err)

	var tries int
	err = b.Do(context.Background(), func(ctx context.Context) error {
		tries++

		return retryableError
	})

	require.Error(t, err)
	assert.Equal(t, maxRetries+1, tries)
}

func TestWithMaxDuration(t *testing.T) {
	const maxDuration = time.Millisecond * 5

	b, err := New(Constant, time.Millisecond, WithMaxDuration(maxDuration))
	require.NoError(t, err)

	started := time.Now()
	err = b.Do(context.Background(), func(ctx context.Context) error {
		return retryableError
	})
	duration := time.Now().Sub(started)

	require.Error(t, err)
	assert.Less(t, duration-maxDuration, time.Millisecond)
}

func TestWithCap(t *testing.T) {
	const cap = time.Millisecond

	b, err := New(Exponential, time.Millisecond, WithCap(cap))
	require.NoError(t, err)

	retries := make([]time.Duration, 0)
	previousStart := time.Now()
	err = b.Do(context.Background(), func(ctx context.Context) error {
		if len(retries) >= 5 {
			return nil
		}

		retries = append(retries, time.Now().Sub(previousStart))

		previousStart = time.Now()

		return retryableError
	})

	assert.NoError(t, err)

	const threshold = cap * 2
	for _, dur := range retries {
		if dur > threshold {
			t.Errorf("time between retries exceeded threshold (cap=%s, threshold=%s, actual=%s", cap, threshold, dur)
		}
	}
}

func TestWithJitter(t *testing.T) {
	const backoff = time.Millisecond * 10
	const jitter = time.Millisecond * 8

	b, err := New(Constant, backoff, WithJitter(jitter))

	retries := make([]time.Duration, 0)
	previousStart := time.Now()
	err = b.Do(context.Background(), func(ctx context.Context) error {
		if len(retries) >= 5 {
			return nil
		}

		retries = append(retries, time.Now().Sub(previousStart))
		previousStart = time.Now()

		return retryableError
	})

	assert.NoError(t, err)

	const threshold = time.Millisecond * 5
	var found bool
	for _, dur := range retries {
		if dur-backoff > threshold || backoff-dur > threshold {
			found = true
			break
		}
	}

	assert.True(t, found, "no jitter found over threshold (=%s)", threshold)
}
