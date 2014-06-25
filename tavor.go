package tavor

import (
	"fmt"
	"io"
	"strings"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

const (
	Version = "0.1"
)

//TODO remove this
var DEBUG = false

func PrettyPrintTree(w io.Writer, root token.Token) {
	prettyPrintTreeRek(w, root, 0)
}

func prettyPrintTreeRek(w io.Writer, tok token.Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case token.ForwardToken:
		if v := t.Get(); v != nil {
			prettyPrintTreeRek(w, v, level+1)
		}
	case lists.List:
		for i := 0; i < t.Len(); i++ {
			c, _ := t.Get(i)

			prettyPrintTreeRek(w, c, level+1)
		}
	}
}

func PrettyPrintInternalTree(w io.Writer, root token.Token) {
	prettyPrintInternalTreeRek(w, root, 0)
}

func prettyPrintInternalTreeRek(w io.Writer, tok token.Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case token.ForwardToken:
		if v := t.InternalGet(); v != nil {
			prettyPrintInternalTreeRek(w, v, level+1)
		}
	case lists.List:
		for i := 0; i < t.InternalLen(); i++ {
			c, _ := t.InternalGet(i)

			prettyPrintInternalTreeRek(w, c, level+1)
		}
	}
}

func LoopExists(root token.Token) bool {
	lookup := make(map[token.Token]struct{})
	queue := linkedlist.New()

	queue.Push(root)

	for !queue.Empty() {
		v, _ := queue.Shift()
		t, _ := v.(token.Token)

		lookup[t] = struct{}{}

		switch tok := t.(type) {
		case *primitives.Pointer:
			if v := tok.InternalGet(); v != nil {
				if _, ok := lookup[v]; ok {
					if DEBUG {
						fmt.Printf("Found a loop through (%p)%+v\n", t)
					}

					return true
				}

				queue.Push(v)
			}
		case token.ForwardToken:
			if v := tok.InternalGet(); v != nil {
				queue.Push(v)
			}
		case lists.List:
			for i := 0; i < tok.InternalLen(); i++ {
				c, _ := tok.InternalGet(i)

				queue.Push(c)
			}
		}
	}

	return false
}
