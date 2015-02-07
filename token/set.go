package token

import (
	"github.com/zimmski/container/list/linkedlist"

	"github.com/zimmski/tavor/log"
)

// ResetResetTokens resets all tokens in the token graph that fullfill the ResetToken interface
func ResetResetTokens(root Token) {
	var queue = linkedlist.New()

	queue.Unshift(root)

	for !queue.Empty() {
		v, _ := queue.Shift()

		switch tok := v.(type) {
		case ResetToken:
			log.Debugf("reset %#v(%p)", tok, tok)

			tok.Reset()
		}

		if t, ok := v.(Follow); ok && !t.Follow() {
			continue
		}

		switch tok := v.(type) {
		case ForwardToken:
			if v := tok.Get(); v != nil {
				queue.Unshift(v)
			}
		case ListToken:
			for i := tok.Len() - 1; i >= 0; i-- {
				c, _ := tok.Get(i)

				queue.Unshift(c)
			}
		}
	}
}

// ResetScope resets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func ResetScope(root Token) {
	log.Debug("start reseting scope")

	SetScope(root, NewVariableScope())

	log.Debug("finished reseting scope")
}

// SetScope sets all scopes of tokens in the token graph that fullfill the ScopeToken interface
func SetScope(root Token, scope *VariableScope) {
	queue := linkedlist.New()

	type set struct {
		token Token
		scope *VariableScope
	}

	queue.Unshift(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		tok := s.token
		scope := s.scope

		if _, ok := tok.(Scoping); ok {
			scope = scope.Push()
		}

		if t, ok := tok.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p) with %#v", t, t, scope)

			t.SetScope(scope)
		}

		if t, ok := tok.(Follow); ok && !t.Follow() {
			continue
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.Get(); v != nil {
				queue.Unshift(set{
					token: v,
					scope: scope,
				})
			}
		case ListToken:
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				queue.Unshift(set{
					token: c,
					scope: scope,
				})
			}
		}
	}
}

// ResetInternalScope resets all scopes of interal tokens in the token graph that fullfill the ScopeToken interface
func ResetInternalScope(root Token) {
	log.Debug("start reseting internal scope")

	SetInternalScope(root, NewVariableScope())

	log.Debug("finished reseting internal scope")
}

// SetInternalScope sets all scopes of internal tokens in the token graph that fullfill the ScopeToken interface
func SetInternalScope(root Token, scope *VariableScope) {
	queue := linkedlist.New()

	type set struct {
		token Token
		scope *VariableScope
	}

	queue.Unshift(set{
		token: root,
		scope: scope,
	})

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		tok := s.token
		scope := s.scope

		if _, ok := tok.(Scoping); ok {
			scope = scope.Push()
		}

		if t, ok := tok.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p) with %#v", t, t, scope)

			t.SetScope(scope)
		}

		if t, ok := tok.(Follow); ok && !t.Follow() {
			continue
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				queue.Unshift(set{
					token: v,
					scope: scope,
				})
			}
		case ListToken:
			for i := t.InternalLen() - 1; i >= 0; i-- {
				c, _ := t.InternalGet(i)

				queue.Unshift(set{
					token: c,
					scope: scope,
				})
			}
		}
	}
}

// ResetCombinedScope resets all scopes of external and interal tokens in the token graph that fullfill the ScopeToken interface
func ResetCombinedScope(root Token) {
	log.Debug("start reseting internal scope")

	SetCombinedScope(root, NewVariableScope())

	log.Debug("finished reseting internal scope")
}

// SetCombinedScope sets all scopes of external and internal tokens in the token graph that fullfill the ScopeToken interface
func SetCombinedScope(root Token, scope *VariableScope) {
	queue := linkedlist.New()
	checked := make(map[Token]struct{})

	type set struct {
		token Token
		scope *VariableScope
	}

	queue.Unshift(set{
		token: root,
		scope: scope,
	})
	checked[root] = struct{}{}

	for !queue.Empty() {
		v, _ := queue.Shift()
		s := v.(set)

		tok := s.token
		scope := s.scope

		if _, ok := tok.(Scoping); ok {
			scope = scope.Push()
		}

		if t, ok := tok.(ScopeToken); ok {
			log.Debugf("setScope %#v(%p) with %#v", t, t, scope)

			t.SetScope(scope)
		}

		if t, ok := tok.(Follow); ok && !t.Follow() {
			continue
		}

		switch t := tok.(type) {
		case ForwardToken:
			if v := t.InternalGet(); v != nil {
				if _, ok := checked[v]; !ok {
					queue.Unshift(set{
						token: v,
						scope: scope,
					})
					checked[v] = struct{}{}
				}
			}
			if v := t.Get(); v != nil {
				if _, ok := checked[v]; !ok {
					queue.Unshift(set{
						token: v,
						scope: scope,
					})
					checked[v] = struct{}{}
				}
			}
		case ListToken:
			for i := t.InternalLen() - 1; i >= 0; i-- {
				c, _ := t.InternalGet(i)

				if _, ok := checked[c]; !ok {
					queue.Unshift(set{
						token: c,
						scope: scope,
					})
					checked[c] = struct{}{}
				}
			}
			for i := t.Len() - 1; i >= 0; i-- {
				c, _ := t.Get(i)

				if _, ok := checked[c]; !ok {
					queue.Unshift(set{
						token: c,
						scope: scope,
					})
					checked[c] = struct{}{}
				}
			}
		}
	}
}
