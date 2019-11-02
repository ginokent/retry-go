package retry

import (
	"errors"
	"time"
)

var (
	errorMaxRetriesExceeded = errors.New("max retries exceeded")
	errorTimeout            = errors.New("timeout")
	errorDoNotReuseRetrier  = errors.New("do not reuse retrier")
)

type retrier struct {
	retries       int
	maxRetries    int
	retryInterval time.Duration
	timeout       time.Duration
	timeoutChan   <-chan time.Time
	error         error
}

// New return retrier
func New(maxRetries int, retryInterval, timeout time.Duration) *retrier {
	return &retrier{
		retries:       0,
		maxRetries:    maxRetries,
		retryInterval: retryInterval,
		timeout:       timeout,
		timeoutChan:   nil,
		error:         nil,
	}
}

func (r *retrier) GetRetries() int            { return r.retries }
func (r *retrier) GetMaxRetries() int         { return r.maxRetries }
func (r *retrier) GetInterval() time.Duration { return r.retryInterval }
func (r *retrier) GetTimeout() time.Duration  { return r.timeout }
func (r *retrier) GetError() error            { return r.error }

func (r *retrier) ResetRetries() {
	r.retries = 0
}

func (r *retrier) ResetError() {
	r.error = error(nil)
}

func (r *retrier) Retry() bool {
	// ** DO NOT REUSE retrier **
	if r.error != nil {
		r.error = errorDoNotReuseRetrier
		return false
	}
	// retry loop
	select {
	case <-r.timeoutChan:
		r.error = errorTimeout
		return false
	default:
		switch {
		case r.retries <= 0:
			// 1 回目は眠らない。まだリトライじゃないから。
			// noop
			r.retries = 1
			r.timeoutChan = time.After(r.timeout)
			return true
		case r.retries > r.maxRetries:
			r.error = errorMaxRetriesExceeded
			return false
		default:
			// 2 回目（リトライ初回）以降は眠る。
			time.Sleep(r.retryInterval)
			r.retries++
			return true
		}
	}
}

func (r *retrier) RetryWithExponentialBackoff() bool {
	// ** DO NOT REUSE retrier **
	if r.error != nil {
		r.error = errorDoNotReuseRetrier
		return false
	}
	// retry loop
	select {
	case <-r.timeoutChan:
		r.error = errorTimeout
		return false
	default:
		switch {
		case r.retries <= 0:
			// 1 回目は眠らない。まだリトライじゃないから。
			// noop
			r.retries = 1
			r.timeoutChan = time.After(r.timeout)
			return true
		case r.retries > r.maxRetries:
			r.error = errorMaxRetriesExceeded
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
