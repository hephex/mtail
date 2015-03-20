package watcher

import (
	"testing"
	"time"
)

func TestFakeWatcher(t *testing.T) {
	w := NewFakeWatcher()

	w.Add("/tmp")
	if _, ok := w.watches["/tmp"]; !ok {
		t.Errorf("Not watching /tmp, w contains: %+#v", w.watches)
	}

	w.Remove("/tmp")
	if _, ok := w.watches["/tmp"]; ok {
		t.Errorf("Still watching /tmp, w contains: %+#v", w.watches)
	}

	w.Add("/tmp")
	w.InjectCreate("/tmp")
	select {
	case name := <-w.Creates():
		if name != "/tmp" {
			t.Errorf("event doesn't match: %s\n", name)
		}
	case <-time.After(1 * time.Millisecond):
		t.Fatalf("No message.")
	}

	w.Add("/tmp/foo")
	w.InjectUpdate("/tmp/foo")
	select {
	case name := <-w.Updates():
		if name != "/tmp/foo" {
			t.Errorf("event doesn't match name: %s\n", name)
		}
	case <-time.After(1 * time.Millisecond):
		t.Fatalf("no event found in watcher: %+#v\n", w)
	}

	w.Close()
}