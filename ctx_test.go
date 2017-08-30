package zerolog

import (
	"context"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestCtx(t *testing.T) {
	log := New(ioutil.Discard)
	ctx := log.WithContext(context.Background())
	log2 := Ctx(ctx)
	if !reflect.DeepEqual(log, *log2) {
		t.Error("Ctx did not return the expected logger")
	}

	// update
	log = log.Level(InfoLevel)
	ctx = log.WithContext(ctx)
	log2 = Ctx(ctx)
	if !reflect.DeepEqual(log, *log2) {
		t.Error("Ctx did not return the expected logger")
	}

	log2 = Ctx(context.Background())
	if log2 != disabledLogger {
		t.Error("Ctx did not return the expected logger")
	}
}

func TestCtxDisabled(t *testing.T) {
	ctx := disabledLogger.WithContext(context.Background())
	if ctx != context.Background() {
		t.Error("WithContext stored a disabled logger")
	}

	ctx = New(ioutil.Discard).WithContext(ctx)
	if reflect.DeepEqual(Ctx(ctx), disabledLogger) {
		t.Error("WithContext did not store logger")
	}

	ctx = disabledLogger.WithContext(ctx)
	if !reflect.DeepEqual(Ctx(ctx), disabledLogger) {
		t.Error("WithContext did not update logger pointer with disabled logger")
	}
}
