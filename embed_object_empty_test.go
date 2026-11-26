package zerolog

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/rs/zerolog/internal/cbor"
)

type emptyObjectMarshaler struct{}

func (emptyObjectMarshaler) MarshalZerologObject(*Event) {}

func TestContextEmbedObjectEmptyProducesValidJSON(t *testing.T) {
	tests := []struct {
		name string
		log  func(buf *bytes.Buffer)
		want string
	}{
		{
			name: "empty marshaler between fields",
			log: func(buf *bytes.Buffer) {
				logger := New(buf).With().
					Str("a", "1").
					EmbedObject(emptyObjectMarshaler{}).
					Str("b", "2").
					Logger()
				logger.Info().Msg("hi")
			},
			want: `{"level":"info","a":"1","b":"2","message":"hi"}` + "\n",
		},
		{
			name: "empty marshaler trailing context",
			log: func(buf *bytes.Buffer) {
				logger := New(buf).With().
					Str("a", "1").
					EmbedObject(emptyObjectMarshaler{}).
					Logger()
				logger.Info().Msg("hi")
			},
			want: `{"level":"info","a":"1","message":"hi"}` + "\n",
		},
		{
			name: "nil marshaler after existing field",
			log: func(buf *bytes.Buffer) {
				logger := New(buf).With().
					Str("a", "1").
					EmbedObject(nil).
					Logger()
				logger.Info().Msg("hi")
			},
			want: `{"level":"info","a":"1","message":"hi"}` + "\n",
		},
		{
			name: "empty marshaler on empty context",
			log: func(buf *bytes.Buffer) {
				logger := New(buf).With().
					EmbedObject(emptyObjectMarshaler{}).
					Str("a", "1").
					Logger()
				logger.Info().Msg("hi")
			},
			want: `{"level":"info","a":"1","message":"hi"}` + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			tc.log(&buf)
			got := cbor.DecodeIfBinaryToBytes(buf.Bytes())
			if !json.Valid(got) {
				t.Fatalf("invalid JSON output: %q", string(got))
			}
			if string(got) != tc.want {
				t.Errorf("got %q, want %q", string(got), tc.want)
			}
		})
	}
}
