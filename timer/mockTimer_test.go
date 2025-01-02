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

	if err := mockTimer.WaitCb(d, cb); err != nil{
		t.Fatal(err)
	}

	if err := mockTimer.ForceDone(); err != nil{
		t.Fatal(err)
	}

	wg.Wait()
}
