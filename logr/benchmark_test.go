package logr_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/logr"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	os.Stderr, _ = os.Open("/dev/null")
}

//go:noinline
func doInfoOneArg(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.Info("this is", "a", "string")
	}
}

//go:noinline
func doInfoSeveralArgs(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.Info("multi",
			"bool", true, "string", "str", "int", 42,
			"float", 3.14, "struct", struct{ X, Y int }{93, 76})
	}
}

//go:noinline
func doV0Info(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.V(0).Info("multi",
			"bool", true, "string", "str", "int", 42,
			"float", 3.14, "struct", struct{ X, Y int }{93, 76})
	}
}

//go:noinline
func doV9Info(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.V(9).Info("multi",
			"bool", true, "string", "str", "int", 42,
			"float", 3.14, "struct", struct{ X, Y int }{93, 76})
	}
}

//go:noinline
func doError(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	err := fmt.Errorf("error message")
	for i := 0; i < b.N; i++ {
		log.Error(err, "multi",
			"bool", true, "string", "str", "int", 42,
			"float", 3.14, "struct", struct{ X, Y int }{93, 76})
	}
}

//go:noinline
func doWithValues(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := log.WithValues("k1", "v1", "k2", "v2")
		_ = l
	}
}

//go:noinline
func doWithName(b *testing.B, log logr.Logger) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		l := log.WithName("name")
		_ = l
	}
}

func logger() logr.Logger {
	zl := zerolog.New(os.Stderr)
	return logr.New(&zl)
}

func BenchmarkDiscardInfoOneArg(b *testing.B) {
	var log logr.Logger = logger()
	doInfoOneArg(b, log)
}

func BenchmarkDiscardInfoSeveralArgs(b *testing.B) {
	var log logr.Logger = logger()
	doInfoSeveralArgs(b, log)
}

func BenchmarkDiscardV0Info(b *testing.B) {
	var log logr.Logger = logger()
	doV0Info(b, log)
}

func BenchmarkDiscardV9Info(b *testing.B) {
	var log logr.Logger = logger()
	doV9Info(b, log)
}

func BenchmarkDiscardError(b *testing.B) {
	var log logr.Logger = logger()
	doError(b, log)
}

func BenchmarkDiscardWithValues(b *testing.B) {
	var log logr.Logger = logger()
	doWithValues(b, log)
}

func BenchmarkDiscardWithName(b *testing.B) {
	var log logr.Logger = logger()
	doWithName(b, log)
}

func BenchmarkFuncrInfoOneArg(b *testing.B) {
	var log logr.Logger = logger()
	doInfoOneArg(b, log)
}

func BenchmarkFuncrInfoSeveralArgs(b *testing.B) {
	var log logr.Logger = logger()
	doInfoSeveralArgs(b, log)
}

func BenchmarkFuncrV0Info(b *testing.B) {
	var log logr.Logger = logger()
	doV0Info(b, log)
}

func BenchmarkFuncrV9Info(b *testing.B) {
	var log logr.Logger = logger()
	doV9Info(b, log)
}

func BenchmarkFuncrError(b *testing.B) {
	var log logr.Logger = logger()
	doError(b, log)
}

func BenchmarkFuncrWithValues(b *testing.B) {
	var log logr.Logger = logger()
	doWithValues(b, log)
}

func BenchmarkFuncrWithName(b *testing.B) {
	var log logr.Logger = logger()
	doWithName(b, log)
}
