package zerolog

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

var logObjectMarshalerType = reflect.TypeOf((*LogObjectMarshaler)(nil)).Elem()

func jsonMarshalNoHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 && b[len(b)-1] == '\n' {
		b = b[:len(b)-1]
	}
	return b, nil
}

func implementsLogObjectMarshaler(t reflect.Type) bool {
	if t == nil {
		return false
	}
	if t.Implements(logObjectMarshalerType) {
		return true
	}
	if t.Kind() != reflect.Ptr {
		return reflect.PointerTo(t).Implements(logObjectMarshalerType)
	}
	return false
}

// logObjectMarshalerIsPromoted reports whether v implements LogObjectMarshaler
// only because an anonymous field does. Interface() should then walk the outer
// struct so sibling fields are not dropped.
func logObjectMarshalerIsPromoted(v interface{}) bool {
	if v == nil {
		return false
	}
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		if reflect.ValueOf(v).IsNil() {
			return false
		}
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || !implementsLogObjectMarshaler(reflect.TypeOf(v)) {
		return false
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && implementsLogObjectMarshaler(f.Type) {
			return true
		}
	}
	return false
}

func marshalInterface(v interface{}) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return marshalValue(reflect.ValueOf(v))
}

func deref(rv reflect.Value) reflect.Value {
	for rv.IsValid() && (rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr) {
		if rv.IsNil() {
			return reflect.Value{}
		}
		rv = rv.Elem()
	}
	return rv
}

func marshalValue(rv reflect.Value) ([]byte, error) {
	if !rv.IsValid() {
		return []byte("null"), nil
	}
	rv = deref(rv)
	if !rv.IsValid() {
		return []byte("null"), nil
	}
	if rv.CanInterface() {
		v := rv.Interface()
		if lom, ok := v.(LogObjectMarshaler); ok && !logObjectMarshalerIsPromoted(v) {
			b := appendObject(nil, lom, false, nil, nil)
			// Object() uses the build-tag encoder. Interface marshaling must
			// stay JSON (including the binary_log embed-JSON path).
			if len(b) > 0 && b[0] == '{' {
				return b, nil
			}
			return jsonMarshalNoHTML(v)
		}
	}
	if rv.Kind() == reflect.Struct && structNeedsLOMWalk(rv.Type()) {
		return marshalStruct(rv)
	}
	if !rv.CanInterface() {
		return []byte("null"), nil
	}
	return jsonMarshalNoHTML(rv.Interface())
}

func structNeedsLOMWalk(t reflect.Type) bool {
	return structNeedsLOMWalkSeen(t, make(map[reflect.Type]bool))
}

func structNeedsLOMWalkSeen(t reflect.Type, seen map[reflect.Type]bool) bool {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return implementsLogObjectMarshaler(t)
	}
	if seen[t] {
		return false
	}
	seen[t] = true
	if implementsLogObjectMarshaler(t) {
		return true
	}
	for i := 0; i < t.NumField(); i++ {
		if structNeedsLOMWalkSeen(t.Field(i).Type, seen) {
			return true
		}
	}
	return false
}

type jsonFieldMeta struct {
	name      string
	skip      bool
	omitempty bool
	inline    bool
}

func jsonFieldInfo(f reflect.StructField) jsonFieldMeta {
	tag := f.Tag.Get("json")
	if tag == "-" {
		return jsonFieldMeta{skip: true}
	}
	name := f.Name
	inline := f.Anonymous
	omitempty := false
	if tag != "" {
		parts := strings.Split(tag, ",")
		if parts[0] != "" {
			name = parts[0]
			inline = false
		}
		for _, p := range parts[1:] {
			if p == "omitempty" {
				omitempty = true
			}
		}
	}
	return jsonFieldMeta{name: name, omitempty: omitempty, inline: inline}
}

func marshalStruct(rv reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	first := true
	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		meta := jsonFieldInfo(f)
		if meta.skip {
			continue
		}
		if f.PkgPath != "" && !f.Anonymous {
			continue
		}
		fv := rv.Field(i)
		if !fv.CanInterface() && !f.Anonymous {
			continue
		}
		if meta.omitempty && isEmptyValue(fv) {
			continue
		}
		b, err := marshalValue(fv)
		if err != nil {
			return nil, err
		}
		if meta.inline {
			inner := bytes.TrimSpace(b)
			if bytes.Equal(inner, []byte("null")) || len(inner) == 0 {
				continue
			}
			if len(inner) >= 2 && inner[0] == '{' && inner[len(inner)-1] == '}' {
				inner = bytes.TrimSpace(inner[1 : len(inner)-1])
				if len(inner) == 0 {
					continue
				}
				if !first {
					buf.WriteByte(',')
				}
				buf.Write(inner)
				first = false
				continue
			}
		}
		if !first {
			buf.WriteByte(',')
		}
		first = false
		key, err := jsonMarshalNoHTML(meta.name)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteByte(':')
		buf.Write(b)
	}
	buf.WriteByte('}')
	return buf.Bytes(), nil
}

func isEmptyValue(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return rv.Len() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return rv.IsNil()
	}
	return false
}
