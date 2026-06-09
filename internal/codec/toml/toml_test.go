package toml

import (
	"reflect"
	"testing"

	"github.com/develonaut/todui/internal/todo"
)

func TestRoundTrip(t *testing.T) {
	in := todo.List{
		Header:      []string{"# T", "", "_sub_"},
		LastUpdated: "2026-06-09 08:50",
		Items: []todo.Item{
			{Title: "a", Section: "now", Order: 0, Tags: []string{"x", "y"}, ADO: "#1"},
			{Title: "b", Section: "next", Order: 0, Description: "ctx"},
			{Title: "d", Section: "done", Order: 0, DoneDate: "2026-06-04"},
		},
	}
	var c Codec
	b, err := c.Encode(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round trip mismatch:\n in=%+v\nout=%+v\n---\n%s", in, out, b)
	}
}

func TestDecodeInvalid(t *testing.T) {
	var c Codec
	if _, err := c.Decode([]byte("this is = not = valid")); err == nil {
		t.Error("expected error decoding invalid TOML")
	}
}
