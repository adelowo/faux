package context_test

import (
	gcontext "context"
	"sync"
	"testing"
	"time"

	"github.com/gu-io/gu/tests"
	"github.com/influx6/faux/context"
)

func TestGoogleContext(t *testing.T) {
	ctx := context.NewGoogleContext(gcontext.Background())
	ctx.Set("brolly", "benzine")

	var wg sync.WaitGroup
	wg.Add(1)

	go goRoutineContext(t, &wg, ctx.New(true))

	go func() {
		<-time.After(1 * time.Millisecond)
		ctx.Cancel()
	}()

	wg.Wait()

	if !ctx.IsExpired() {
		tests.Failed(t, "Should have successfully expired context")
	}
	tests.Passed(t, "Should have successfully expired context")

	_, hasTime := ctx.TimeRemaining()
	if hasTime {
		tests.Failed(t, "Should have time allocated to context")
	}
	tests.Passed(t, "Should have time allocated to context")
}

// TestContextWithConnectedChild tests the validaty of the context.
func TestContextWithConnectedChild(t *testing.T) {
	ctx := context.New()
	ctx.Set("brolly", "benzine")

	var wg sync.WaitGroup
	wg.Add(1)

	go goRoutineContext(t, &wg, ctx.New(true))

	go func() {
		<-time.After(1 * time.Millisecond)
		ctx.Cancel()
	}()

	wg.Wait()

	if !ctx.IsExpired() {
		tests.Failed(t, "Should have successfully expired context")
	}
	tests.Passed(t, "Should have successfully expired context")

	rem, hasTime := ctx.TimeRemaining()
	if hasTime {
		tests.Failed(t, "Should have time allocated to context")
	}
	tests.Passed(t, "Should have time allocated to context")

	if rem != 0 {
		tests.Failed(t, "Should have successfully used up time")
	}
	tests.Passed(t, "Should have successfully used up time")

}

// TestContextWithDisconnectedChild tests the validaty of the context.
func TestContextWithDisconnectedChild(t *testing.T) {
	ctx := context.New()
	ctx.Set("brolly", "benzine")

	var wg sync.WaitGroup
	wg.Add(1)

	go goRoutineContextFailure(t, &wg, ctx.WithTimeout(3*time.Millisecond, false))

	go func() {
		<-time.After(1 * time.Millisecond)
		ctx.Cancel()
	}()

	wg.Wait()
}

// goRoutineContext tests the usage and canceling of the provided context variable.
func goRoutineContext(t *testing.T, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	val, ok := ctx.Get("brolly")
	if !ok {
		tests.Failed(t, "Should have found key %q in the context", "brolly")
	}
	tests.Passed(t, "Should have found key %q in the context", "brolly")

	realVal, _ := val.(string)
	if realVal != "benzine" {
		tests.Failed(t, "Should have matched context key %q with value %q", "brolly", "benzine")
	}
	tests.Passed(t, "Should have matched context key %q with value %q", "brolly", "benzine")

	select {
	case <-time.After(3 * time.Millisecond):
		tests.Failed(t, "Should have successfully died through parent's call")
	case <-ctx.Done():
		tests.Passed(t, "Should have successfully died through parent's call")
	}
}

// goRoutineContextFailure tests the usage and canceling of the provided context variable.
func goRoutineContextFailure(t *testing.T, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	val, ok := ctx.Get("brolly")
	if !ok {
		tests.Failed(t, "Should have found key %q in the context", "brolly")
	}
	tests.Passed(t, "Should have found key %q in the context", "brolly")

	realVal, _ := val.(string)
	if realVal != "benzine" {
		tests.Failed(t, "Should have matched context key %q with value %q", "brolly", "benzine")
	}
	tests.Passed(t, "Should have matched context key %q with value %q", "brolly", "benzine")

	select {
	case <-time.After(2 * time.Millisecond):
		tests.Passed(t, "Should have successfully died through parent's call")
	case <-ctx.Done():
		tests.Failed(t, "Should have successfully died through parent's call")
	}
}
