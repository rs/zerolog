//go:build !binary_log
// +build !binary_log

package zerolog

import (
	"bytes"
	"strings"
	"testing"
)

// Types for #600: MarshalZerologObject should apply when a value is nested
// inside another struct logged via Interface.

type InnerSecret struct {
	Secret string `json:"secret"`
	OK     string `json:"ok"`
}

func (i InnerSecret) MarshalZerologObject(e *Event) {
	e.Str("ok", i.OK)
}

type EmbeddedSecret struct {
	InnerSecret
	Name string `json:"name"`
}

type NestedSecret struct {
	Child InnerSecret `json:"child"`
	Name  string      `json:"name"`
}

func TestInterfaceCallsMarshalZerologObjectOnNestedStructs(t *testing.T) {
	inner := InnerSecret{Secret: "password", OK: "visible"}

	t.Run("direct", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Interface("v", inner).Msg("")
		got := decodeIfBinaryToString(out.Bytes())
		t.Logf("direct: %s", got)
		if strings.Contains(got, "password") {
			t.Fatalf("direct Interface leaked secret: %s", got)
		}
		if !strings.Contains(got, `"ok":"visible"`) {
			t.Fatalf("direct Interface missing redacted fields: %s", got)
		}
	})

	t.Run("embedded", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Interface("v", EmbeddedSecret{InnerSecret: inner, Name: "bob"}).Msg("")
		got := decodeIfBinaryToString(out.Bytes())
		t.Logf("embedded: %s", got)
		if strings.Contains(got, "password") {
			t.Fatalf("embedded Interface leaked secret: %s", got)
		}
		if !strings.Contains(got, `"ok":"visible"`) {
			t.Fatalf("embedded Interface missing redacted fields: %s", got)
		}
		if !strings.Contains(got, `"name":"bob"`) {
			t.Fatalf("embedded Interface dropped outer fields: %s", got)
		}
	})

	t.Run("nested", func(t *testing.T) {
		out := &bytes.Buffer{}
		log := New(out)
		log.Info().Interface("v", NestedSecret{Child: inner, Name: "bob"}).Msg("")
		got := decodeIfBinaryToString(out.Bytes())
		t.Logf("nested: %s", got)
		if strings.Contains(got, "password") {
			t.Fatalf("nested Interface leaked secret: %s", got)
		}
		if !strings.Contains(got, `"ok":"visible"`) {
			t.Fatalf("nested Interface missing redacted fields: %s", got)
		}
	})
}
