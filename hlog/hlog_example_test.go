package hlog_test

import (
	"net/http"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

// fake alice to avoid dep
type alice struct{}

func (a alice) Append(interface{}) alice       { return a }
func (alice) Then(h http.Handler) http.Handler { return h }

func Example_handler() {
	host, _ := os.Hostname()
	log := zerolog.New(os.Stdout).With().
		Str("role", "my-service").
		Str("host", host).
		Logger()

	c := alice{}

	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(log))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to those handler, all our logs will come with some pre-populated fields.
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))

	// Here is your final handler
	h := c.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the logger from the request's context. You can safely assume it
		// will be always there: if the handler is removed, hlog.FromRequest
		// will return a no-op logger.
		hlog.FromRequest(r).Info().
			Str("user", "current user").
			Str("status", "ok").
			Msg("Something happend")
	}))
	http.Handle("/", h)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal().Err(err).Msg("Startup failed")
	}
}
