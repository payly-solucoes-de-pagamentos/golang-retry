package policies_test

import (
	"errors"
	"testing"
	"time"

	policies "github.com/payly-solucoes-de-pagamentos/golang-retry"
	"github.com/stretchr/testify/assert"
)

func TestRetryShouldBeCalledThreeTimes(testing *testing.T) {
	// arrange
	oneSecond := 1 * time.Second
	attempts := uint(3)
	policy := policies.NewRetryPolicy()
	policy.SetAttempts(attempts).SetDelay(oneSecond)
	counter := 0

	// act
	errors := policy.Retry(func() error {
		counter++
		if counter < 3 {
			return errors.New("some retry error")
		}
		return nil
	})

	// assert
	assert.Equal(testing, 3, counter)
	assert.NotNil(testing, errors)
	assert.Len(testing, errors, 2)
}

func TestRetryShouldFailAfterThreeAttempts(testing *testing.T) {
	// arrange
	oneSecond := 1 * time.Second
	attempts := uint(3)
	policy := policies.NewRetryPolicy()
	policy.SetAttempts(attempts).SetDelay(oneSecond)

	// act
	errors := policy.Retry(func() error {
		return errors.New("some retry error")
	})

	// assert
	assert.NotNil(testing, errors)
	assert.Len(testing, errors, 3)
	// assert.Equal(testing, "All attempts fail:\n#1: some retry error\n#2: some retry error\n#3: some retry error", err.Error())
}

func TestRetryShouldBeSucceededAfterSecondAttempt(testing *testing.T) {
	// arrange
	oneSecond := 1 * time.Second
	attempts := uint(3)
	policy := policies.NewRetryPolicy()
	policy.SetAttempts(attempts).SetDelay(oneSecond)
	counter := 0

	// act
	errors := policy.Retry(func() error {
		counter++
		if counter < 2 {
			return errors.New("some retry error")
		}
		return nil
	})

	// assert
	assert.Equal(testing, 2, counter)
	assert.NotNil(testing, errors)
	assert.Len(testing, errors, 1)
}
