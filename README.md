# retry-go

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
