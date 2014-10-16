package graph

import (
	"bytes"
	"testing"

	. "github.com/zimmski/tavor/test/assert"

	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

func TestGraphDot(t *testing.T) {
	// TODO test WriteDot
	// - FIXME that there are no pointer strings in the output -> these change all the time
	// - FIXME order of definitions in the file -> sorting
	/*
		expected := `digraph Graphing {
			node [shape = doublecircle]; c2080365a8;
			node [shape = point] START;
			node [shape = ellipse];

			c208036a10 [label="1"]
			c208036ae0 [label="2"]
			c208036468 [label="3"]
			c2080365a8 [label="4"]

			START -> c208036a10;
			c208036a10 -> c208036ae0;
			c208036a10 -> c208036468;
			c208036ae0 -> c2080365a8;
			c208036468 -> c2080365a8;
		}
		`
	*/

	var got bytes.Buffer

	WriteDot(lists.NewAll(
		primitives.NewConstantInt(1),
		lists.NewOne(
			primitives.NewConstantInt(2),
			primitives.NewConstantInt(3),
		),
		primitives.NewConstantInt(4),
	), &got)

	True(t, len(got.String()) > 0)
}
