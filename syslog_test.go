// +build !binary_log
// +build !windows

package zerolog

import "testing"
import "reflect"

type syslogEvent struct {
	level string
	msg   string
}
type syslogTestWriter struct {
	events []syslogEvent
}

func (w *syslogTestWriter) Write(p []byte) (int, error) {
	return 0, nil
}
func (w *syslogTestWriter) Debug(m string) error {
	w.events = append(w.events, syslogEvent{"Debug", m})
	return nil
}
func (w *syslogTestWriter) Info(m string) error {
	w.events = append(w.events, syslogEvent{"Info", m})
	return nil
}
func (w *syslogTestWriter) Warning(m string) error {
	w.events = append(w.events, syslogEvent{"Warning", m})
	return nil
}
func (w *syslogTestWriter) Err(m string) error {
	w.events = append(w.events, syslogEvent{"Err", m})
	return nil
}
func (w *syslogTestWriter) Emerg(m string) error {
	w.events = append(w.events, syslogEvent{"Emerg", m})
	return nil
}
func (w *syslogTestWriter) Crit(m string) error {
	w.events = append(w.events, syslogEvent{"Crit", m})
	return nil
}

func TestSyslogWriter(t *testing.T) {
	sw := &syslogTestWriter{}
	log := New(SyslogLevelWriter(sw))
	log.Debug().Msg("debug")
	log.Info().Msg("info")
	log.Warn().Msg("warn")
	log.Error().Msg("error")
	log.Log().Msg("nolevel")
	want := []syslogEvent{
		{"Debug", `{"level":"debug","message":"debug"}` + "\n"},
		{"Info", `{"level":"info","message":"info"}` + "\n"},
		{"Warning", `{"level":"warn","message":"warn"}` + "\n"},
		{"Err", `{"level":"error","message":"error"}` + "\n"},
		{"Info", `{"message":"nolevel"}` + "\n"},
	}
	if got := sw.events; !reflect.DeepEqual(got, want) {
		t.Errorf("Invalid syslog message routing: want %v, got %v", want, got)
	}
}
