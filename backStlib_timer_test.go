//go:build !wasm

package time

import (
	"testing"
	"time"

	"github.com/tinywasm/time"
)

func TestAfterFunc(t *testing.T) {
	tp := time.NewTimeProvider()
	executed := false

	tp.AfterFunc(50, func() {
		executed = true
	})

	if executed {
		t.Error("callback executed too early")
	}

	time.Sleep(100 * time.Millisecond)

	if !executed {
		t.Error("callback was not executed")
	}
	t.Log("AfterFunc test passed")
}

func TestAfterFunc_Stop(t *testing.T) {
	tp := time.NewTimeProvider()
	timer, executed := AfterFuncStopSetup(tp)

	wasActive := timer.Stop()
	if !wasActive {
		t.Error("timer should have been active")
	}

	time.Sleep(200 * time.Millisecond)
	AfterFuncStopVerify(t, executed)
}
