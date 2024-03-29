# retry-go
[![Build Status](https://travis-ci.com/djeeno/retry-go.svg?branch=master)](https://travis-ci.com/djeeno/retry-go)
[![codecov](https://codecov.io/gh/djeeno/retry-go/branch/master/graph/badge.svg)](https://codecov.io/gh/djeeno/retry-go)

## retry
```go
r := retry.New(maxRetries, retryInterval, timeout)

for r.Retry() {
	// do something
	if success {
		break
	}
}

if r.GetError() != nil {
	return r.GetError() // "max retries exceeded" or "timeout"
}
```

## retry with exponential backoff
```go
r := retry.New(maxRetries, retryInterval, timeout)

for r.RetryWithExponentialBackoff() {
	// do something
	if success {
		break
	}
}

if r.GetError() != nil {
	return r.GetError() // "max retries exceeded" or "timeout"
}
```
