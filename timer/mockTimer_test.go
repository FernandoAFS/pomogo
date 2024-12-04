package timer

import (
	"sync"
	"testing"
	"time"
)

// ANY TIME WILL DO...
var refTime = time.Date(2024, 12, 04, 0, 0, 0, 0, time.UTC)

// TEST IF THE CALLBACK IS RUN CORRECTLY WHEN IT SHOULD...
func TestMockTimerSync(t *testing.T) {

	d, err := time.ParseDuration("24h")

	if err != nil {
		t.Fatal(err)
	}

	mockTimer := MockTimerFactory(refTime)

	var wg sync.WaitGroup
	wg.Add(1)

	cb := func(then time.Time) {
		if then != refTime.Add(d) {
			t.Fatalf("Then is not expected value")
		}
		wg.Done()
	}
	mockTimer.WaitCb(d, cb)

	mockTimer.ForceDone()
	wg.Wait()
}
