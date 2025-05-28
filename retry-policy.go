package policies

import (
	"fmt"
	"time"

	"github.com/avast/retry-go/v4"
)

type OnRetryFunc func(attempt uint, err error)

type ExecuteFunc retry.RetryableFunc

type IRetryPolicy interface {
	Retry(execute ExecuteFunc) RetryErrors
	SetAttempts(attempts uint) IRetryPolicy
	SetDelay(delay time.Duration) IRetryPolicy
}

type RetryPolicy struct {
	attempts uint
	delay    time.Duration
	onRetry  OnRetryFunc
}

func (policy *RetryPolicy) SetAttempts(attempts uint) IRetryPolicy {
	policy.attempts = attempts
	return policy
}

func (policy *RetryPolicy) SetDelay(delay time.Duration) IRetryPolicy {
	policy.delay = delay
	return policy
}

func (policy *RetryPolicy) OnRetry(onRetry OnRetryFunc) {
	policy.onRetry = onRetry
}

func (policy *RetryPolicy) addRetryError(retryErrors *RetryErrors) func(attempt uint, err error) {
	return func(attempt uint, err error) {
		retryErr := RetryError{
			Attempt: attempt,
			Message: err.Error(),
		}
		*retryErrors = append(*retryErrors, retryErr)
	}
}

func (policy RetryPolicy) Retry(execute ExecuteFunc) RetryErrors {
	retryErrors := make(RetryErrors, 0)

	exec := retry.RetryableFunc(execute)

	retry.Do(exec,
		retry.OnRetry(policy.addRetryError(&retryErrors)),
		retry.Attempts(policy.attempts),
		retry.Delay(policy.delay),
	)

	return retryErrors
}

func NewRetryPolicy() IRetryPolicy {
	return &RetryPolicy{
		attempts: 3,
		delay:    5 * time.Second,
	}
}

type RetryError struct {
	Attempt uint   `json:"attempt"`
	Message string `json:"error"`
}

type RetryErrors []RetryError

func (e RetryError) Error() string {
	return fmt.Sprintf("Attempt %d got error message %s", e.Attempt, e.Message)
}

func (errs RetryErrors) ToErrorInterface() []error {
	var errors []error
	for _, err := range errs {
		errors = append(errors, err)
	}
	return errors
}
