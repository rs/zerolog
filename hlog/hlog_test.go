// +build go1.7

package hlog

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"reflect"

	"net/http/httptest"

	"github.com/rs/zerolog"
)

func TestNewHandler(t *testing.T) {
	log := zerolog.New(nil).With().
		Str("foo", "bar").
		Logger()
	lh := NewHandler(log)
	h := lh(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		if !reflect.DeepEqual(l, log) {
			t.Fail()
		}
	}))
	h.ServeHTTP(nil, &http.Request{})
}

func TestURLHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		URL: &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h := URLHandler("url")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"url":"/path?foo=bar"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestMethodHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Method: "POST",
	}
	h := MethodHandler("method")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"method":"POST"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestRequestHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/path", RawQuery: "foo=bar"},
	}
	h := RequestHandler("request")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"request":"POST /path?foo=bar"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestRemoteAddrHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		RemoteAddr: "1.2.3.4:1234",
	}
	h := RemoteAddrHandler("ip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"ip":"1.2.3.4"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestRemoteAddrHandlerIPv6(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		RemoteAddr: "[2001:db8:a0b:12f0::1]:1234",
	}
	h := RemoteAddrHandler("ip")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"ip":"2001:db8:a0b:12f0::1"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestUserAgentHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"User-Agent": []string{"some user agent string"},
		},
	}
	h := UserAgentHandler("ua")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"ua":"some user agent string"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestRefererHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"Referer": []string{"http://foo.com/bar"},
		},
	}
	h := RefererHandler("referer")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := `{"referer":"http://foo.com/bar"}`+"\n", out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(nil, r)
}

func TestRequestIDHandler(t *testing.T) {
	out := &bytes.Buffer{}
	r := &http.Request{
		Header: http.Header{
			"Referer": []string{"http://foo.com/bar"},
		},
	}
	h := RequestIDHandler("id", "Request-Id")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := IDFromRequest(r)
		if !ok {
			t.Fatal("Missing id in request")
		}
		if want, got := id.String(), w.Header().Get("Request-Id"); got != want {
			t.Errorf("Invalid Request-Id header, got: %s, want: %s", got, want)
		}
		l := FromRequest(r)
		l.Log().Msg("")
		if want, got := fmt.Sprintf(`{"id":"%s"}`+"\n", id), out.String(); want != got {
			t.Errorf("Invalid log output, got: %s, want: %s", got, want)
		}
	}))
	h = NewHandler(zerolog.New(out))(h)
	h.ServeHTTP(httptest.NewRecorder(), r)
}
