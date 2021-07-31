// +build !binary_log

package logr_test

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/logr"
)

type E struct {
	str string
}

func (e E) Error() string {
	return e.str
}

func Helper(log logr.Logger, msg string) {
	helper2(log, msg)
}

func helper2(log logr.Logger, msg string) {
	log.WithCallDepth(2).Info(msg)
}

func ExampleNew() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	zl := zerolog.New(os.Stdout)
	log := logr.New(&zl)
	log = log.WithName("MyName")
	log = log.WithValues("module", "example")

	log.Info("hello", "val1", 1, "val2", map[string]int{"k": 1})
	log.V(1).Info("you should see this")
	log.V(1).V(1).Info("you should NOT see this")
	log.Error(nil, "uh oh", "trouble", true, "reasons", []float64{0.1, 0.11, 3.14})
	log.Error(E{"an error occurred"}, "goodbye", "code", -1)
	Helper(log, "thru a helper")

	// Output:
	// {"level":"info","module":"example","val1":1,"val2":{"k":1},"logger":"MyName","message":"hello"}
	// {"level":"debug","module":"example","logger":"MyName","message":"you should see this"}
	// {"level":"error","module":"example","trouble":true,"reasons":[0.1,0.11,3.14],"logger":"MyName","message":"uh oh"}
	// {"level":"error","error":"an error occurred","module":"example","code":-1,"logger":"MyName","message":"goodbye"}
	// {"level":"info","message":"thru a helper"}
}
