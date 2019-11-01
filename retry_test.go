package retry

import (
	"log"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	testRetrier := New(1, 100*time.Millisecond, 10*time.Second)
	testRetrier.GetRetries()
	testRetrier.GetMaxRetries()
	testRetrier.GetInterval()
	testRetrier.GetTimeout()
	want := error(nil)
	actual := testRetrier.Error()
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
		if errorMaxRetriesExceeded != testMaxRetriesExceeded.Error() {
			t.Fatalf("want: %v, actual: %v", errorMaxRetriesExceeded, testMaxRetriesExceeded.Error())
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
		if errorTimeout != testTimeout.Error() {
			t.Fatalf("want: %v, actual: %v", errorTimeout, testTimeout.Error())
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
		if errorDoNotReuseRetrier != testReuse.Error() {
			t.Fatalf("want: %v, actual: %v", errorDoNotReuseRetrier, testReuse.Error())
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
		if errorMaxRetriesExceeded != testMaxRetriesExceeded.Error() {
			t.Fatalf("want: %v, actual: %v", errorMaxRetriesExceeded, testMaxRetriesExceeded.Error())
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
		if errorTimeout != testTimeout.Error() {
			t.Fatalf("want: %v, actual: %v", errorTimeout, testTimeout.Error())
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
		if errorDoNotReuseRetrier != testReuse.Error() {
			t.Fatalf("want: %v, actual: %v", errorDoNotReuseRetrier, testReuse.Error())
		}
	}
}
