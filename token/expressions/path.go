package expressions

import (
	"bytes"
	"github.com/zimmski/tavor/token/primitives"

	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/token"
	"github.com/zimmski/tavor/token/lists"
	"github.com/zimmski/tavor/token/variables"
)

// Path implements a path query
type Path struct {
	list      token.Token
	from      token.Token
	over      token.Token
	connectBy []token.Token
	without   []token.Token

	variableScope *token.VariableScope
}

// NewPath returns a new instance of a Path token given the set of tokens
func NewPath(list token.Token, from token.Token, over token.Token, connectBy []token.Token, without []token.Token) (*Path, error) {
	if err := checkListToken(list); err != nil {
		return nil, err
	}

	return &Path{
		list:      list,
		from:      from,
		over:      over,
		connectBy: connectBy,
		without:   without,
	}, nil
}

func checkListToken(list token.Token) error {
	if token.LoopExists(list) {
		return &token.ParserError{
			Message: "There is an endless loop in the list argument. Use a variable reference to avoid this.",
			Type:    token.ParseErrorEndlessLoopDetected,
		}
	}

	return nil
}

func (e *Path) path() []string {
	variableScope := e.variableScope.Push()

	if p, ok := e.list.(*primitives.Pointer); ok {
		e.list = p.Resolve()
	}

	tl := e.list

	if v, ok := tl.(*primitives.Scope); ok {
		tl = v.Get()
	}

	if t, ok := tl.(token.ScopeToken); ok {
		t.SetScope(variableScope)
	}

	if v, ok := tl.(*variables.VariableReference); ok {
		tl = v.Reference()
	}

	if v, ok := tl.(*primitives.Scope); ok {
		tl = v.Get()
	}

	l, ok := tl.(token.ListToken)
	if !ok {
		// TODO must be a ListToken but is ...

		return nil
	}

	connects := make(map[string][]string, 0)

	for i := 0; i < l.Len(); i++ {
		el, _ := l.Get(i)
		if v, ok := el.(*primitives.Scope); ok {
			el = v.Get()
		}

		variableScope.Set("e", variables.NewVariable("e", el))

		if t, ok := e.over.(token.ScopeToken); ok {
			token.SetScope(t, variableScope)
		}
		for j := 0; j < len(e.connectBy); j++ {
			if t, ok := e.connectBy[j].(token.ScopeToken); ok {
				token.SetScope(t, variableScope)
			}
		}

		cs := make([]string, len(e.connectBy))
		for j := 0; j < len(e.connectBy); j++ {
			cs[j] = e.connectBy[j].String()
		}

		connects[e.over.String()] = cs
	}

	token.SetScope(e.from, variableScope)
	from := e.from.String()

	path := []string{from}

	checked := make(map[string]struct{})
	checked[from] = struct{}{}
	for i := 0; i < len(e.without); i++ {
		checked[e.without[i].String()] = struct{}{}
	}

	stack := linkedlist.New()
	stack.Unshift(from)

	for !stack.Empty() {
		v, _ := stack.Shift()
		c := v.(string)

		n, ok := connects[c]
		if !ok {
			continue
		}

		for i := len(n) - 1; i >= 0; i-- {
			v := n[i]

			if _, ok := checked[v]; !ok {
				path = append(path, v)

				checked[v] = struct{}{}
				stack.Unshift(v)
			}
		}
	}

	return path
}

// Token interface methods

// Clone returns a copy of the token and all its children
func (e *Path) Clone() token.Token {
	return &Path{
		list:      e.list,
		from:      e.from,
		over:      e.over,
		connectBy: e.connectBy,
		without:   e.without,

		variableScope: e.variableScope,
	}
}

// Parse tries to parse the token beginning from the current position in the parser data.
// If the parsing is successful the error argument is nil and the next current position after the token is returned.
func (e *Path) Parse(pars *token.InternalParser, cur int) (int, []error) {
	panic("TODO")
}

// Permutation sets a specific permutation for this token
func (e *Path) Permutation(i uint) error {
	permutations := e.Permutations()

	if i < 1 || i > permutations {
		return &token.PermutationError{
			Type: token.PermutationErrorIndexOutOfBound,
		}
	}

	// do nothing

	return nil
}

// Permutations returns the number of permutations for this token
func (e *Path) Permutations() uint {
	return 1
}

// PermutationsAll returns the number of all possible permutations for this token including its children
func (e *Path) PermutationsAll() uint {
	return 1
}

func (e *Path) String() string {
	var buffer bytes.Buffer

	for _, el := range e.path() {
		if _, err := buffer.WriteString(el); err != nil {
			panic(err)
		}
	}

	return buffer.String()
}

// List interface methods

// Get returns the current referenced token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *Path) Get(i int) (token.Token, error) {
	l := e.path()

	if i < 0 || i >= len(l) {
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}

	return primitives.NewConstantString(l[i]), nil
}

// Len returns the number of the current referenced tokens
func (e *Path) Len() int {
	return len(e.path())
}

// InternalGet returns the current referenced internal token at the given index. The error return argument is not nil, if the index is out of bound.
func (e *Path) InternalGet(i int) (token.Token, error) {
	il := e.InternalLen()

	if i < 0 || i >= il {
		return nil, &lists.ListError{
			Type: lists.ListErrorOutOfBound,
		}
	}

	if i == 0 {
		return e.list, nil
	} else if i == 1 {
		return e.over, nil
	} else if i < 2+len(e.connectBy) {
		return e.connectBy[i-2], nil
	}

	return e.without[i-2-len(e.connectBy)], nil
}

// InternalLen returns the number of referenced internal tokens
func (e *Path) InternalLen() int {
	return 2 + len(e.connectBy) + len(e.without)
}

// InternalLogicalRemove removes the referenced internal token and returns the replacement for the current token or nil if the current token should be removed.
func (e *Path) InternalLogicalRemove(tok token.Token) token.Token {
	panic("TODO")
}

// InternalReplace replaces an old with a new internal token if it is referenced by this token. The error return argument is not nil, if the replacement is not suitable.
func (e *Path) InternalReplace(oldToken, newToken token.Token) error {
	if e.list == oldToken {
		if err := checkListToken(newToken); err != nil {
			return err
		}

		e.list = newToken
	}

	if e.over == oldToken {
		e.over = newToken
	}

	for i := 0; i < len(e.connectBy); i++ {
		if e.connectBy[i] == oldToken {
			e.connectBy[i] = newToken
		}
	}

	for i := 0; i < len(e.without); i++ {
		if e.without[i] == oldToken {
			e.without[i] = newToken
		}
	}

	return nil
}

// ScopeToken interface methods

var _ token.ScopeToken = (*Path)(nil)

// SetScope sets the scope of the token
func (e *Path) SetScope(variableScope *token.VariableScope) {
	e.variableScope = variableScope
}
