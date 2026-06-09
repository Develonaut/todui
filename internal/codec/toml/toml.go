// Package toml implements ports.Codec using TOML as the on-disk format.
package toml

import (
	gotoml "github.com/pelletier/go-toml/v2"

	"github.com/develonaut/todui/internal/ports"
	"github.com/develonaut/todui/internal/todo"
)

// Codec encodes and decodes a todo.List as TOML. Items are emitted as an
// array of tables ([[item]]); short arrays such as tags stay inline.
type Codec struct{}

// Encode marshals the list to TOML bytes.
func (Codec) Encode(l todo.List) ([]byte, error) {
	return gotoml.Marshal(l)
}

// Decode unmarshals TOML bytes into a list.
func (Codec) Decode(b []byte) (todo.List, error) {
	var l todo.List
	if err := gotoml.Unmarshal(b, &l); err != nil {
		return todo.List{}, err
	}
	return l, nil
}

var _ ports.Codec = Codec{}
