package time

import (
	"testing"

	"github.com/tinywasm/time"
)

// AfterFuncStopShared tests that Stop() prevents callback execution.
// The waiting mechanism is platform-specific, so this only sets up the test.
// Returns: timer, pointer to executed flag
func AfterFuncStopSetup(tp time.TimeProvider) (time.Timer, *bool) {
	executed := false
	timer := tp.AfterFunc(100, func() {
		executed = true
	})
	return timer, &executed
}

// AfterFuncStopVerify checks that the callback was NOT executed after Stop()
func AfterFuncStopVerify(t *testing.T, executed *bool) {
	if *executed {
		t.Error("callback should not have executed after Stop()")
	}
	t.Log("AfterFunc_Stop test passed")
}
