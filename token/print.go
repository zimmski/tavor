package token

import (
	"fmt"
	"io"
	"strings"
)

// PrettyPrintTree prints the token represenation as a text tree
func PrettyPrintTree(w io.Writer, root Token) {
	prettyPrintTreeRek(w, root, 0)
}

func prettyPrintTreeRek(w io.Writer, tok Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case ForwardToken:
		if v := t.Get(); v != nil {
			prettyPrintTreeRek(w, v, level+1)
		}
	case ListToken:
		for i := 0; i < t.Len(); i++ {
			c, _ := t.Get(i)

			prettyPrintTreeRek(w, c, level+1)
		}
	}
}

// PrettyPrintInternalTree prints the internal token represenation as a text tree
func PrettyPrintInternalTree(w io.Writer, root Token) {
	prettyPrintInternalTreeRek(w, root, 0)
}

func prettyPrintInternalTreeRek(w io.Writer, tok Token, level int) {
	fmt.Fprintf(w, "%s(%p)%#v\n", strings.Repeat("\t", level), tok, tok)

	switch t := tok.(type) {
	case ForwardToken:
		if v := t.InternalGet(); v != nil {
			prettyPrintInternalTreeRek(w, v, level+1)
		}
	case ListToken:
		for i := 0; i < t.InternalLen(); i++ {
			c, _ := t.InternalGet(i)

			prettyPrintInternalTreeRek(w, c, level+1)
		}
	}
}
