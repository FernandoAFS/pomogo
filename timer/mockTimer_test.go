package timer

import (
	"sync"
	"testing"
	"time"
)

// TEST IF THE CALLBACK IS RUN CORRECTLY WHEN IT SHOULD...
func TestMockTimerSync(t *testing.T) {

	d, err := time.ParseDuration("24h")

	if err != nil {
		t.Fatal(err)
	}

	mockTimer := MockCbTimer{}

	var wg sync.WaitGroup
	wg.Add(1)

	cb := func() {
		wg.Done()
	}
	mockTimer.WaitCb(d, cb)

	mockTimer.ForceDone()
	wg.Wait()
}
