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
	log2, ok := FromContext(ctx)
	if !ok {
		t.Error("Expected ok=true from FromContext")
	}
	if !reflect.DeepEqual(log, log2) {
		t.Error("FromContext did not return the expected logger")
	}

	log2, ok = FromContext(context.Background())
	if ok {
		t.Error("Expected ok=false from FromContext")
	}
	if !reflect.DeepEqual(log2, Logger{}) {
		t.Error("FromContext did not return the expected logger")
	}
}
