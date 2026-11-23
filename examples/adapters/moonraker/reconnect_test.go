package moonraker_test

import (
	"testing"
	"time"

	"example.com/prusa/moonraker/moonrakersim"
)

func TestKlippyDisconnectEndsSessionForRedial(t *testing.T) {
	t.Parallel()

	machine := startMachine(t, moonrakersim.Options{})
	session := connect(t, machine)
	collect(t, session, map[string]string{"/state": `"IDLE"`})

	machine.Notify("notify_klippy_disconnected")
	timer := time.NewTimer(settle)
	defer timer.Stop()
	for {
		select {
		case _, open := <-session.State():
			if !open {
				machine.Push(map[string]map[string]any{"print_stats": {"state": "printing", "filename": "test.gcode"}})
				reconnected := connect(t, machine)
				collect(t, reconnected, map[string]string{"/state": `"PRINTING"`})
				machine.Push(map[string]map[string]any{"print_stats": {"state": "paused"}})
				collect(t, reconnected, map[string]string{"/state": `"PAUSED"`})
				return
			}
		case <-timer.C:
			t.Fatal("session remained open after Klippy disconnected")
		}
	}
}
