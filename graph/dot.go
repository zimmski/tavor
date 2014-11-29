package graph

import (
	"fmt"
	"io"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/constraints"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/primitives"
)

type dotGraph struct {
	edges    []dotEdge
	vertices map[token.Token]dotVertice

	end token.Token
}

type dotEdge struct {
	from, to token.Token
	optional bool
	label    string
}

type dotVertice struct {
	label string
	typ   string
}

func nodeUID(tok token.Token) string {
	return fmt.Sprintf("%p", tok)[1:]
}

func (g *dotGraph) addDot(tok token.Token) (start, next map[token.Token]bool) {
	start = make(map[token.Token]bool)
	next = make(map[token.Token]bool)

	switch t := tok.(type) {
	case *constraints.Optional:
		v := t.InternalGet()

		s, n := g.addDot(v)

		for k := range s {
			start[k] = true
		}
		for k := range n {
			next[k] = true
		}

		start[g.end] = false
		next[g.end] = false
	case *lists.One:
		l := t.InternalLen()

		for i := 0; i < l; i++ {
			c, _ := t.InternalGet(i)

			if o, ok := c.(token.OptionalToken); ok && o.IsOptional() {
				start[g.end] = false
				next[g.end] = false
			}

			s, n := g.addDot(c)

			for k, v := range s {
				start[k] = v
			}
			for k, v := range n {
				next[k] = v
			}
		}
	case *lists.All:
		var prev map[token.Token]bool

		l := t.InternalLen()

		for i := 0; i < l; i++ {
			c, _ := t.InternalGet(i)

			s, n := g.addDot(c)

			if i == 0 {
				for k, v := range s {
					if k == g.end {
						continue
					}

					start[k] = v
				}
			} else {
				found := false

				for pt, po := range prev {
					for st, so := range s {
						if pt == g.end || st == g.end {
							continue
						}

						g.edges = append(g.edges, dotEdge{
							from:     pt,
							to:       st,
							optional: po || so,
						})

						if so {
							found = true
						}
					}
				}

				if found {
					for pt := range prev {
						if pt == g.end {
							continue
						}

						n[pt] = false
					}
				}
			}

			prev = n

			if i == l-1 || l == 1 {
				for k, v := range n {
					if k == g.end {
						continue
					}

					next[k] = v
				}
			}
		}
	case *lists.Repeat:
		var times string
		if t.From() == t.To() {
			times = fmt.Sprintf("%dx", t.From())
		} else {
			times = fmt.Sprintf("%d-%dx", t.From(), t.To())
		}

		a := primitives.NewConstantString("")
		g.vertices[a] = dotVertice{
			label: a.String(),
			typ:   "point",
		}

		b := primitives.NewConstantString("")
		g.vertices[b] = dotVertice{
			label: b.String(),
			typ:   "point",
		}

		v, _ := t.InternalGet(0)

		s, n := g.addDot(v)

		for st, so := range s {
			g.edges = append(g.edges, dotEdge{
				from:     a,
				to:       st,
				optional: so,
			})
		}

		for nt, no := range n {
			g.edges = append(g.edges, dotEdge{
				from:     nt,
				to:       b,
				optional: no,
			})
		}

		g.edges = append(g.edges, dotEdge{
			from:     b,
			to:       a,
			optional: false,
			label:    times,
		})

		if t.From() == 0 {
			start[a] = true
			next[b] = true
		} else {
			start[a] = false
			next[b] = false
		}
	case token.ListToken:
		panic(fmt.Errorf("%#v not implemented", t))
	case *primitives.RangeInt:
		label := fmt.Sprintf("%d..%d", t.From(), t.To())

		if st := t.Step(); st != 1 {
			label += fmt.Sprintf(" with step %d", st)
		}

		g.vertices[tok] = dotVertice{
			label: label,
		}

		start[tok] = false
		next[tok] = false
	default:
		g.vertices[tok] = dotVertice{
			label: tok.String(),
		}

		start[tok] = false
		next[tok] = false
	}

	return
}

// WriteDot writes the token graph as DOT format to the writer
func WriteDot(root token.Token, dst io.Writer) {
	g := &dotGraph{
		vertices: make(map[token.Token]dotVertice),

		end: primitives.NewConstantString("END"),
	}

	start, next := g.addDot(root)

	foundEnd := false

	for k := range start {
		if k == g.end {
			foundEnd = true

			break
		}
	}
	for _, edge := range g.edges {
		if edge.to == g.end {
			foundEnd = true

			break
		}
	}

	if foundEnd {
		g.vertices[g.end] = dotVertice{
			label: g.end.String(),
		}
	}

	fmt.Fprintf(dst, "digraph Graphing {\n")

	fmt.Fprintf(dst, "\tnode [shape = doublecircle];")
	for tok := range next {
		fmt.Fprintf(dst, " %s", nodeUID(tok))
	}
	fmt.Fprintf(dst, ";\n")

	fmt.Fprintf(dst, "\tnode [shape = point] START;\n")

	for tok, vertice := range g.vertices {
		if vertice.typ != "" {
			fmt.Fprintf(dst, "\tnode [shape = %s] %s;\n", vertice.typ, nodeUID(tok))
		}
	}

	fmt.Fprintf(dst, "\tnode [shape = ellipse];\n")

	fmt.Fprintln(dst)

	for tok, vertice := range g.vertices {
		fmt.Fprintf(dst, "\t%s [label=%q]\n", nodeUID(tok), vertice.label)
	}

	fmt.Fprintln(dst)

	found := false

	for tok, opt := range start {
		fmt.Fprintf(dst, "\tSTART -> %s", nodeUID(tok))

		if opt {
			fmt.Fprintf(dst, " [style=dotted]")

			found = true
		}

		fmt.Fprintf(dst, ";\n")
	}

	if found {
		for _, edge := range g.edges {
			if edge.optional {
				fmt.Fprintf(dst, "\tSTART -> %s;\n", nodeUID(edge.to))
			}
		}
	}

	for _, edge := range g.edges {
		fmt.Fprintf(dst, "\t%s -> %s", nodeUID(edge.from), nodeUID(edge.to))

		if edge.optional || edge.label != "" {
			fmt.Fprint(dst, "[")

			if edge.optional {
				fmt.Fprintf(dst, " style=dotted")
			}
			if edge.label != "" {
				fmt.Fprintf(dst, " label=%q", edge.label)
			}

			fmt.Fprint(dst, "]")
		}

		fmt.Fprintf(dst, ";\n")
	}

	fmt.Fprintf(dst, "}\n")

	/*
		digraph graphname {

			size="8,5"
			node [shape = doublecircle]; LR_0 LR_3 LR_4 LR_8;
			node [shape = circle];

		    a [label="Google"]
		    b [label="Apple"]
		    c [label="UCLA"]
		    d [label="Stanford"]
		    a -> b [label=50, color=blue];
		    b -> c [label=-10, color=red];
		    b -> d [label="A", color=green];
			b -- d [style=dotted];
		}
	*/
}
