package journald

import "testing"

func TestWriteWithoutJournal(t *testing.T) {
	w := NewJournalDWriter()
	n, err := w.Write([]byte(`{"level":"info","message":"test"}`))
	if err == nil {
		t.Fatal("expected an error when journald is unavailable on Windows")
	}
	if n != 0 {
		t.Fatalf("expected no bytes written, got %d", n)
	}
}
