// Package jsoncodec installs encoding/json/v2 as Gin's JSON implementation.
// Request binding and response rendering both go through this codec.
package jsoncodec

import (
	"bytes"
	"errors"
	"io"

	jsonv1 "encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	ginjson "github.com/gin-gonic/gin/codec/json"
)

// Install replaces Gin's JSON codec with encoding/json/v2.
// v2 defaults apply: case-sensitive names, duplicate names rejected,
// invalid UTF-8 rejected, and nil slices and maps encoded as empty
// arrays and objects. HTML escaping stays off unless the encoder's
// SetEscapeHTML(true) is called.
func Install() {
	ginjson.API = api{}
}

type api struct{}

func (api) Marshal(v any) ([]byte, error) {
	return jsonv2.Marshal(v)
}

func (api) Unmarshal(data []byte, v any) error {
	return jsonv2.Unmarshal(data, v)
}

func (api) MarshalIndent(v any, prefix, indent string) ([]byte, error) {
	if prefix == "" && indent == "" {
		return jsonv2.Marshal(v)
	}
	return jsonv2.Marshal(v, jsontext.WithIndentPrefix(prefix), jsontext.WithIndent(indent))
}

func (api) NewEncoder(w io.Writer) ginjson.Encoder {
	return &encoder{w: w}
}

func (api) NewDecoder(r io.Reader) ginjson.Decoder {
	return &decoder{dec: jsontext.NewDecoder(r)}
}

type encoder struct {
	w          io.Writer
	escapeHTML bool
}

func (e *encoder) SetEscapeHTML(on bool) {
	e.escapeHTML = on
}

func (e *encoder) Encode(v any) error {
	var opts []jsonv2.Options
	if e.escapeHTML {
		opts = append(opts, jsontext.EscapeForHTML(true))
	}
	if err := jsonv2.MarshalWrite(e.w, v, opts...); err != nil {
		return err
	}
	_, err := e.w.Write([]byte{'\n'})
	return err
}

type decoder struct {
	dec           *jsontext.Decoder
	useNumber     bool
	rejectUnknown bool
}

func (d *decoder) UseNumber() {
	d.useNumber = true
}

func (d *decoder) DisallowUnknownFields() {
	d.rejectUnknown = true
}

func (d *decoder) Decode(v any) error {
	opts := make([]jsonv2.Options, 0, 2)
	if d.useNumber {
		opts = append(opts, jsonv2.WithUnmarshalers(useNumberUnmarshaler))
	}
	if d.rejectUnknown {
		opts = append(opts, jsonv2.RejectUnknownMembers(true))
	}
	return jsonv2.UnmarshalDecode(d.dec, v, opts...)
}

// useNumberUnmarshaler matches encoding/json.Decoder.UseNumber: JSON numbers
// stored in an any become encoding/json.Number. Every other kind uses v2's
// default any representation.
var useNumberUnmarshaler = jsonv2.UnmarshalFromFunc(func(dec *jsontext.Decoder, v *any) error {
	if dec.PeekKind() != jsontext.KindNumber {
		return errors.ErrUnsupported
	}
	raw, err := dec.ReadValue()
	if err != nil {
		return err
	}
	*v = jsonv1.Number(bytes.TrimSpace(raw))
	return nil
})
