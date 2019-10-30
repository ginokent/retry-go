package retry

import (
	"errors"
	"time"
)

var (
	ErrorMaxRetriesExceeded = errors.New("maxRetries retries exceeded")
	ErrorTimeout            = errors.New("timeout")
)

type Retrier struct {
	retries       int
	maxRetries    int
	retryInterval time.Duration
	timeout       <-chan time.Time
	error         error
}

func New(maxRetries int, retryInterval, timeout time.Duration) *Retrier {
	return &Retrier{
		retries:       0,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		timeout:       time.After(timeout),
		error:         nil,
	}
}

func (r *Retrier) GetRetries() int              { return r.retries }
func (r *Retrier) GetMaxRetries() int           { return r.maxRetries }
func (r *Retrier) GetInterval() time.Duration   { return r.retryInterval }
func (r *Retrier) GetTimeout() <-chan time.Time { return r.timeout }
func (r *Retrier) Error() error                 { return r.error }

func (r *Retrier) Retry() bool {
	select {
	case <-r.timeout:
		r.error = ErrorTimeout
		return false
	default:
		switch {
		case r.retries == 0:
			// 1 回目は眠らない。まだリトライじゃないから。
			// noop
			r.retries++
			return true
		case r.retries > r.maxRetries:
			r.error = ErrorMaxRetriesExceeded
			return false
		default:
			// 2 回目（リトライ初回）以降は眠る。
			time.Sleep(r.retryInterval)
			r.retries++
			return true
		}
	}
}

func (r *Retrier) RetryWithExponentialBackoff() bool {
	select {
	case <-r.timeout:
		r.error = ErrorTimeout
		return false
	default:
		switch {
		case r.retries == 0:
			// 1 回目は眠らない。まだリトライじゃないから。
			// noop
			r.retries++
			return true
		case r.retries > r.maxRetries:
			r.error = ErrorMaxRetriesExceeded
			return false
		default:
			// 2 回目（リトライ初回）以降は眠る。
			// 2 回目は retryInterval*1 ( retryInterval<<(retries-1) == retryInterval<<0 == retryInterval*1 ) 眠る
			// 3 回目は retryInterval*2 ( retryInterval<<(retries-1) == retryInterval<<1 == retryInterval*2 ) 眠る
			// 4 回目は retryInterval*4 ( retryInterval<<(retries-1) == retryInterval<<2 == retryInterval*4 ) 眠る
			// 5 回目は retryInterval*8 ( retryInterval<<(retries-1) == retryInterval<<3 == retryInterval*8 ) 眠る
			time.Sleep(time.Duration(r.retryInterval << (r.retries - 1)))
			r.retries++
			return true
		}
	}
}
