package retry

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	maxRetries := 1
	retryInterval := 100 * time.Millisecond
	timeout := 100 * time.Millisecond
	testRetrier := New(maxRetries, retryInterval, timeout)
	if testRetrier.GetRetries() != 0 {
		t.Errorf("TestNew(): testRetrier.GetRetries() != 0")
	}
	if testRetrier.GetMaxRetries() != maxRetries {
		t.Errorf("TestNew(): testRetrier.GetMaxRetries() != maxRetries")
	}
	if testRetrier.GetInterval() != retryInterval {
		t.Errorf("TestNew(): testRetrier.GetInterval() != retryInterval")
	}
	if testRetrier.GetTimeout() != timeout {
		t.Errorf("TestNew(): testRetrier.GetTimeout() != timeout")
	}
	if testRetrier.ResetRetries(); testRetrier.GetRetries() != 0 {
		t.Errorf("TestNew(): testRetrier.ResetRetries(); testRetrier.GetRetries() != 0")
	}
	testRetrier.error = errors.New("test error")
	testRetrier.ResetError()
	want := error(nil)
	actual := testRetrier.GetError()
	if want != actual {
		t.Errorf("want: %v, actual: %v", want, actual)
	}
}

func TestRetrier_Retry(t *testing.T) {
	// errorMaxRetriesExceeded
	{
		testMaxRetriesExceeded := New(4, 100*time.Millisecond, 10*time.Second)
		log.Printf("start: %#v\n", testMaxRetriesExceeded)
		for testMaxRetriesExceeded.Retry() {
			log.Printf("  tmp: %#v\n", testMaxRetriesExceeded)
		}
		log.Printf("  end: %#v\n", testMaxRetriesExceeded)
		if errorMaxRetriesExceeded != testMaxRetriesExceeded.GetError() {
			t.Errorf("want: %v, actual: %v", errorMaxRetriesExceeded, testMaxRetriesExceeded.GetError())
		}
	}

	// errorTimeout
	{
		testTimeout := New(4, 100*time.Millisecond, 100*time.Millisecond)
		log.Printf("start: %#v\n", testTimeout)
		for testTimeout.Retry() {
			log.Printf("  tmp: %#v\n", testTimeout)
		}
		log.Printf("  end: %#v\n", testTimeout)
		if errorTimeout != testTimeout.GetError() {
			t.Errorf("want: %v, actual: %v", errorTimeout, testTimeout.GetError())
		}
	}

	// errorDoNotReuseRetrier
	{
		testReuse := New(4, 100*time.Millisecond, 100*time.Millisecond)
		log.Printf("start: %#v\n", testReuse)
		for testReuse.Retry() {
		}
		log.Printf("  mid: %#v\n", testReuse)
		for testReuse.Retry() {
		}
		log.Printf("  end: %#v\n", testReuse)
		if errorDoNotReuseRetrier != testReuse.GetError() {
			t.Errorf("want: %v, actual: %v", errorDoNotReuseRetrier, testReuse.GetError())
		}
	}
}

func TestRetry_NewSleepExponentialBackoff(t *testing.T) {
	// errorMaxRetriesExceeded
	{
		var before, after time.Time
		before = time.Now()
		testMaxRetriesExceeded := New(4, 100*time.Millisecond, 10*time.Second)
		log.Printf("start: %#v\n", testMaxRetriesExceeded)
		for testMaxRetriesExceeded.RetryWithExponentialBackoff() {
			log.Printf("  tmp: %#v\n", testMaxRetriesExceeded)
			after = time.Now()
			delta := after.Sub(before).Milliseconds()
			want := strconv.Itoa(int(delta))
			actual := strconv.Itoa((1 << (testMaxRetriesExceeded.GetRetries() - 1)) / 2)
			if ! strings.HasPrefix(want, actual) {
				t.Errorf("sleep time is wrong. want: %v, actual: %v", want, actual)
			}
			before = after
		}
		log.Printf("  end: %#v\n", testMaxRetriesExceeded)
		if errorMaxRetriesExceeded != testMaxRetriesExceeded.GetError() {
			t.Errorf("want: %v, actual: %v", errorMaxRetriesExceeded, testMaxRetriesExceeded.GetError())
		}
	}

	// errorTimeout
	{
		var before, after time.Time
		before = time.Now()
		testTimeout := New(4, 100*time.Millisecond, 1*time.Second)
		log.Printf("start: %#v\n", testTimeout)
		for testTimeout.RetryWithExponentialBackoff() {
			log.Printf("  tmp: %#v\n", testTimeout)
			after = time.Now()
			delta := after.Sub(before).Milliseconds()
			want := strconv.Itoa(int(delta))
			actual := strconv.Itoa((1 << (testTimeout.GetRetries() - 1)) / 2)
			if ! strings.HasPrefix(want, actual) {
				t.Errorf("sleep time is wrong. want: %v, actual: %v", want, actual)
			}
			before = after
		}
		log.Printf("  end: %#v\n", testTimeout)
		if errorTimeout != testTimeout.GetError() {
			t.Errorf("want: %v, actual: %v", errorTimeout, testTimeout.GetError())
		}
	}

	// errorDoNotReuseRetrier
	{
		testReuse := New(4, 100*time.Millisecond, 100*time.Millisecond)
		log.Printf("start: %#v\n", testReuse)
		for testReuse.RetryWithExponentialBackoff() {
		}
		log.Printf("  mid: %#v\n", testReuse)
		for testReuse.RetryWithExponentialBackoff() {
		}
		log.Printf("  end: %#v\n", testReuse)
		if errorDoNotReuseRetrier != testReuse.GetError() {
			t.Errorf("want: %v, actual: %v", errorDoNotReuseRetrier, testReuse.GetError())
		}
	}
}
