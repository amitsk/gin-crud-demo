package jsoncodec

import (
	"bytes"
	"strings"
	"testing"

	jsonv1 "encoding/json"
	jsonv2 "encoding/json/v2"

	ginjson "github.com/gin-gonic/gin/codec/json"
)

func TestInstallUsesV2Semantics(t *testing.T) {
	Install()

	type payload struct {
		Names []string `json:"names"`
	}

	out, err := ginjson.API.Marshal(payload{})
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"names":[]}` {
		t.Fatalf("nil slice encoded as %s, want empty array", out)
	}

	err = ginjson.API.Unmarshal([]byte(`{"names":"a","names":"b"}`), &payload{})
	if err == nil {
		t.Fatal("duplicate object names were accepted")
	}
}

func TestEncoderEscapeHTML(t *testing.T) {
	Install()

	var plain bytes.Buffer
	enc := ginjson.API.NewEncoder(&plain)
	if err := enc.Encode(map[string]string{"html": "<tag>"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(plain.String(), "<tag>") {
		t.Fatalf("v2 encoder escaped HTML by default: %s", plain.String())
	}
	if !strings.HasSuffix(plain.String(), "\n") {
		t.Fatalf("encoder output missing trailing newline: %q", plain.String())
	}

	var escaped bytes.Buffer
	enc = ginjson.API.NewEncoder(&escaped)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(map[string]string{"html": "<tag>"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(escaped.String(), `\u003c`) {
		t.Fatalf("SetEscapeHTML(true) left raw HTML: %s", escaped.String())
	}
}

func TestDecoderOptions(t *testing.T) {
	Install()

	dec := ginjson.API.NewDecoder(strings.NewReader(`{"extra":1,"n":2}`))
	dec.UseNumber()
	dec.DisallowUnknownFields()

	var dst struct {
		N any `json:"n"`
	}
	err := dec.Decode(&dst)
	if err == nil {
		t.Fatal("unknown member was accepted")
	}

	dec = ginjson.API.NewDecoder(strings.NewReader(`{"n":2}`))
	dec.UseNumber()
	dst = struct {
		N any `json:"n"`
	}{}
	if err := dec.Decode(&dst); err != nil {
		t.Fatal(err)
	}
	num, ok := dst.N.(jsonv1.Number)
	if !ok || num.String() != "2" {
		t.Fatalf("number decoded as %#v, want json.Number(\"2\")", dst.N)
	}
}

func TestMarshalIndent(t *testing.T) {
	Install()

	out, err := ginjson.API.MarshalIndent(map[string]int{"a": 1}, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	want := "{\n  \"a\": 1\n}"
	if string(out) != want {
		t.Fatalf("indent = %s, want %s", out, want)
	}

	// Compact form stays identical to Marshal.
	compact, err := ginjson.API.MarshalIndent(map[string]int{"a": 1}, "", "")
	if err != nil {
		t.Fatal(err)
	}
	plain, err := jsonv2.Marshal(map[string]int{"a": 1})
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(compact, plain) {
		t.Fatalf("empty indent = %s, marshal = %s", compact, plain)
	}
}
