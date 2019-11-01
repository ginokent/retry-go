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

if r.Error() != nil {
	return r.Error() // "max retries exceeded" or "timeout"
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

if r.Error() != nil {
	return r.Error() // "max retries exceeded" or "timeout"
}
```
