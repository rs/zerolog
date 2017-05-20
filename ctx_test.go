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
	if !reflect.DeepEqual(log, log2) {
		t.Error("Ctx did not return the expected logger")
	}

	// update
	log = log.Level(InfoLevel)
	ctx = log.WithContext(ctx)
	log2 = Ctx(ctx)
	if !reflect.DeepEqual(log, log2) {
		t.Error("Ctx did not return the expected logger")
	}

	log2 = Ctx(context.Background())
	if !reflect.DeepEqual(log2, disabledLogger) {
		t.Error("Ctx did not return the expected logger")
	}
}
